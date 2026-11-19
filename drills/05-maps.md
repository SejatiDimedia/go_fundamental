# Catatan Deep Dive: Maps di Go

Dokumen ini merangkum seluruh hasil diskusi dan deep dive pada sub-topik **2. Composite Types > Maps**, mencakup model memori `hmap`, inisialisasi, comma-ok idiom, syarat map key, hingga isu fatal konkurensi di production.

---

## 1. Anatomi & Struktur Memori Map (`runtime.hmap`)

Di balik layar, tipe `map[K]V` di Go bukanlah struct biasa seperti slice header, melainkan **sebuah pointer langsung ke struct `runtime.hmap`**.

```go
// Representasi konseptual runtime.hmap di Go source code
type hmap struct {
    count     int         // Jumlah elemen dalam map (len)
    flags     uint8       // Status operasi (misal: hashWriting flag saat ada goroutine menulis)
    B         uint8       // Logaritma basis 2 dari jumlah bucket (2^B buckets)
    hash0     uint32      // Hash seed acak per instance map
    buckets   unsafe.Pointer // Pointer ke array bucket (bmap)
    // ... field lainnya
}
```

* **Buckets (`bmap`)**: Setiap bucket menampung hingga **8 pasangan key-value**.
* Jika bucket penuh, Go mengalokasikan **overflow bucket**.
* Jika load factor terlalu tinggi, Go melakukan proses *evacuation / rehashing* ke bucket array yang 2x lebih besar secara bertahap (*incremental evacuation*).

---

## 2. Inisialisasi: `make()` vs `var m map` (The Nil Map Trap)

| Operasi | Nil Map (`var m map[string]int`) | Initialized Map (`m := make(map[string]int)`) |
| :--- | :--- | :--- |
| **Pointer State** | `m == nil` (Pointer `hmap` kosong) | Pointer menunjuk ke `hmap` valid di heap |
| **Read (`v := m["k"]`)** | **Aman** (mengembalikan zero value) | Aman (mengembalikan value atau zero value) |
| **Delete (`delete(m, "k")`)**| **Aman** (no-op, tidak error) | Aman (menghapus jika ada) |
| **Len (`len(m)`)** | **Aman** (mengembalikan `0`) | Aman (mengembalikan jumlah elemen) |
| **Write (`m["k"] = 10`)** | ❌ **PANIC:** `assignment to entry in nil map` | ✅ **Berhasil** |

> [!WARNING]
> Jangan pernah menulis ke nil map. Selalu inisialisasi map dengan `make(map[K]V)` atau map literal `map[K]V{}` sebelum melakukan operasi assignment!

---

## 3. Syarat Tipe Data Map Key: *Comparable Types*

Key sebuah map **wajib mendukung operator perbandingan `==` dan `!=`**:

* ✅ **Boleh jadi Key:**
  - Tipe primitif: `int`, `float64`, `string`, `bool`
  - Pointer: `*T`
  - Array berukuran tetap yang elemennya comparable: `[3]int`, `[2]string`
  - Struct yang semua field-nya comparable
  - Channel: `chan T`
  - Interface (jika dynamic valuenya comparable)

* ❌ **DILARANG jadi Key (Compile Error: `invalid map key type`):**
  - **Slice (`[]T`)**
  - **Map (`map[K]V`)**
  - **Function (`func()`)**

---

## 4. Comma-ok Idiom: Membedakan Missing Key vs Zero Value

Di Go, membaca key yang tidak ada **tidak akan melempar error atau exception**, melainkan mengembalikan zero value dari tipe value-nya.

```go
scores := map[string]int{
    "Bob": 0, // Bob memang bernilai 0
}

// ❌ SALAH: Mengecek dengan membandingkan zero value
if scores["Alice"] == 0 {
    // Ambigu! Apakah Alice tidak ada di map, atau Alice bernilai 0?
}

// ✅ BENAR: Gunakan Comma-ok Idiom
val, ok := scores["Alice"]
if !ok {
    fmt.Println("Alice tidak terdaftar")
} else {
    fmt.Println("Skor Alice:", val)
}
```

---

## 5. Iterasi Map: Sengaja Tidak Deterministik (Randomized Order)

Ketika melakukan `for k, v := range m`, urutan elemen yang keluar selalu acak / berbeda antar eksekusi program.
* **Alasan Go Team:** Go sengaja menambahkan offset acak (*hash seed*) pada awal iterasi untuk mencegah developer bergantung pada urutan map. Jika butuh urutan terurut, kumpulkan key ke dalam slice lalu sort secara eksplisit menggunakan `sort.Strings()` atau `slices.Sort()`.

---

## 6. Concurrency & Race Condition: The Fatal Production Crash

### Kenapa Map Tidak Thread-Safe secara Default?
1. **Zero-Overhead Principle**: Mayoritas map di aplikasi Go digunakan dalam lingkup single-goroutine. Menambahkan mutex internal pada setiap operasi map akan memperlambat performa 90%+ kasus penggunaan yang tidak membutuhkan konkurensi.
2. **Explicit is Better than Implicit**: Go menuntut developer secara sadar memilih mekanisme sinkronisasi yang sesuai dengan workload.

### Akibat Fatal Jika Terjadi Concurrent Write:
Ketika 2 goroutine mengakses map yang sama secara bersamaan (satu menulis, satu membaca/menulis):
* Go Runtime mendeteksi benturan flag (`flags & hashWriting`).
* Go Runtime memicu **`fatal error: concurrent map read and map write`**.
* **CRITICAL PITFALL:** Fatal error ini **TIDAK BISA di-`recover()`** dengan `recover()`! Aplikasi / container pods Kubernetes akan langsung crash instan (*process exit status 2*).

### Solusi Standar Industri:
1. **`sync.RWMutex` (Paling sering digunakan):**
   ```go
   type SafeCounter struct {
       mu   sync.RWMutex
       data map[string]int
   }
   
   func (c *SafeCounter) Inc(key string) {
       c.mu.Lock()
       defer c.mu.Unlock()
       c.data[key]++
   }
   
   func (c *SafeCounter) Get(key string) (int, bool) {
       c.mu.RLock()
       defer c.mu.RUnlock()
       v, ok := c.data[key]
       return v, ok
   }
   ```
2. **`sync.Map`:** Khusus untuk skenario di mana key stabil (sedikit penambahan) dan membaca berkali-kali lipat lebih sering daripada menulis (*cache-like pattern*).

---

## 7. The Addressability Trap: Struct sebagai Value Map

```go
type User struct { Score int }
m := map[string]User{"budi": {Score: 10}}
m["budi"].Score = 20 // ❌ Compile Error: cannot assign to struct field in map
```

### Mengapa Dilarang?
1. Operasi akses `m["budi"]` mengembalikan **temporary copy (salinan sementara)** dari value, bukan pointer ke memori di bucket.
2. Elemen map **tidak addressable** (kamu tidak bisa mengambil alamatnya dengan `&m["budi"]`), karena ukuran map dinamis dan bucket dapat berpindah alamat memori di heap saat terjadi rehashing/evacuation.

### 2 Cara Memperbaiki:
* **Solusi 1 (Pointer Value):** Gunakan `map[string]*User`.
  ```go
  m := map[string]*User{"budi": &User{Score: 10}}
  m["budi"].Score = 20 // ✅ Berhasil karena value adalah pointer
  ```
* **Solusi 2 (Copy -> Mutate -> Re-assign):**
  ```go
  u := m["budi"]
  u.Score = 20
  m["budi"] = u // ✅ Simpan balik ke map
  ```

---

## 8. Map Memory Retention (Leak): `delete()` Tidak Mengecilkan RAM

Di Go, fungsi `delete(m, key)` **TIDAK PERNAH membebaskan bucket memori** yang sudah dialokasikan oleh `runtime.hmap`.

```text
Isi 500.000 elemen    -> Alokasi RAM: 144 MB
delete() semua elemen  -> Alokasi RAM: 144 MB (TETAP SAMA!)
Re-assign map baru     -> Alokasi RAM: 0 MB   (Memori bebas!)
```

* `delete()` hanya mengosongkan entri slot di dalam bucket, namun jumlah bucket tidak pernah di-*shrink*.
* **Solusi Production:** Untuk membebaskan RAM dari map cache yang sempat membengkak, re-assign variabel map ke map baru (`m = make(...)`) atau `nil` agar map raksasa yang lama dapat disapu bersih oleh Garbage Collector.
