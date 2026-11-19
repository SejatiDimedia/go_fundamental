# 📖 Buku Saku: Karakteristik, Larangan Keras & Kaidah Emas Go (Topik 1 - 6)

Dokumen ini adalah **panduan rujukan visual cepat** yang merangkum sifat memori, larangan keras (pantangan compiler/runtime), serta kaidah idiomatik Go **lengkap dengan contoh kode perbandingan (`❌ SALAH` vs `✅ BENAR`)**.

---

## 🗂️ DAFTAR ISI
1. [Modul 1: Language Basics](#modul-1-language-basics)
2. [Modul 2: Composite Types (Array, Slice, Map, Struct)](#modul-2-composite-types)
3. [Modul 3: Control Flow (Conditionals & Loops)](#modul-3-control-flow)
4. [Modul 4: Functions & Closures](#modul-4-functions--closures)
5. [Modul 5: Pointers & Memory Management](#modul-5-pointers--memory-management)
6. [Modul 6: Methods & Interfaces](#modul-6-methods--interfaces)

---

## MODUL 1: Language Basics

### 💡 1. Karakteristik: Zero Value Guarantee
Semua variabel di Go yang tidak diinisialisasi nilai awalnya akan otomatis bersih terisi **Zero Value** (tidak ada memori sampah acak).

```go
var a int     // Otomatis 0
var s string  // Otomatis "" (string kosong)
var b bool    // Otomatis false
var p *int    // Otomatis nil
```

### 💡 2. Karakteristik: Strongly & Statically Typed
Go melarang konversi tipe implisit/otomatis. Semua konversi tipe harus ditulis terang-terangan.

```go
var a int32 = 10
var b int64 = 20

// ❌ COMPILE ERROR: invalid operation: a + b (mismatched types int32 and int64)
// total := a + b

// ✅ BENAR: Konversi tipe eksplisit
total := int64(a) + b
```

### 🚫 Larangan 1: Variabel Lokal Tidak Digunakan
Go melarang adanya variabel lokal yang dideklarasikan tetapi tidak pernah dibaca.

```go
func main() {
    // ❌ COMPILE ERROR: x declared and not used
    x := 10 
    
    // ✅ BENAR: Gunakan variabelnya atau buang dengan blank identifier '_'
    y := 20
    _ = y
}
```

### 🚫 Larangan 2: Mengubah Nilai Konstanta (`const`)
Konstanta ditentukan saat compile-time dan bersifat permanen di memori read-only.

```go
const TaxRate = 0.11

// ❌ COMPILE ERROR: cannot assign to TaxRate (neither addressable nor a map index expression)
// TaxRate = 0.12
```

### 🚫 Larangan 3: Variable Shadowing yang Tidak Disengaja
Menggunakan `:=` di dalam blok inner-scope akan menutupi variabel luar, bukan memodifikasinya!

```go
name := "Alice"
if true {
    // ❌ JEBAKAN: Ini variabel baru di scope if, variabel luar TETAP "Alice"!
    name := "Bob"
    _ = name
}
fmt.Println(name) // Output: tetap "Alice"!

// ✅ BENAR: Gunakan operator '=' biasa tanpa titik dua
if true {
    name = "Bob"
}
fmt.Println(name) // Output: "Bob"
```

---

## MODUL 2: Composite Types

### 1. Arrays & Slices

#### 💡 Karakteristik: Slice Header di RAM
Sebuah slice adalah struct kecil 24-byte di RAM:
* `Data`: pointer ke array asli di RAM
* `Len`: panjang aktif
* `Cap`: kapasitas maksimum array

#### 🚫 Larangan 1: Membiarkan `append` Menimpa Data Slice Lain
Jika dua slice berbagi underlying array yang sama dan kapasitasnya masih ada, `append` pada slice pertama akan menimpa elemen milik slice kedua!

```go
base := []int{1, 2, 3, 4, 5}
s1 := base[0:2] // len: 2, cap: 5 ([1, 2])

// ❌ JEBAKAN: Menimpa base[2] yang tadinya angka 3 menjadi 99!
s1 = append(s1, 99)
fmt.Println(base) // Output: [1 2 99 4 5] -> Rusak!

// ✅ BENAR: Gunakan 3-index slicing [low:high:max] untuk membatasi kapasitas
s2 := base[0:2:2] // len: 2, cap: 2
s2 = append(s2, 99) // Go terpaksa alokasi array baru, base AMAN!
```

#### 🚫 Larangan 2: Memory Leak pada Sub-Slice dari Array Raksasa
Memotong 2 elemen dari file/array 1 GB akan menahan seluruh 1 GB itu di RAM!

```go
hugeArray := make([]byte, 1_000_000_000) // 1 GB

// ❌ JEBAKAN: 1 GB tidak bisa dibersihkan GC karena 'small' menunjuk ke sana!
small := hugeArray[:2]

// ✅ BENAR: Gandakan ke slice baru dengan copy() agar 1 GB bisa dibuang oleh GC
safeSmall := make([]byte, 2)
copy(safeSmall, hugeArray[:2])
```

#### 🚫 Larangan 3: Mengira Slice Biasa Bisa Dimodifikasi Panjangnya oleh Fungsi
Slice Header di-copy by value ke parameter fungsi.

```go
// ❌ GAGAL: Hanya mengubah Len pada fotokopi lokal
func Add(s []int) { s = append(s, 99) }

// ✅ BENAR: Gunakan pointer to slice (*[]T)
func AddPtr(s *[]int) { *s = append(*s, 99) }
```

---

### 2. Maps

#### 💡 Karakteristik: Iterasi yang Diacak & Comma-Ok Idiom
Urutan `for range` map diacak oleh Go. Jika key tidak ada, map mengembalikan zero value tanpa error.

```go
m := map[string]int{"apel": 10}

// ✅ Selalu gunakan Comma-Ok idiom untuk membedakan antara angka 0 vs key tidak ada:
val, ok := m["jeruk"]
if !ok {
    fmt.Println("Key tidak ditemukan!")
}
```

#### 🚫 Larangan 1: DILARANG Menulis ke Nil Map
Membaca dari nil map aman, tetapi MENULIS ke nil map memicu panic seketika!

```go
var m map[string]int // Nil map (belum di-make)

// ❌ RUNTIME PANIC: assignment to entry in nil map
// m["admin"] = 100 

// ✅ BENAR: Selalu inisialisasi dulu
m = make(map[string]int)
m["admin"] = 100
```

#### 🚫 Larangan 2: Key Map Wajib Bertipe `comparable`
Slice, Map, dan Function dilarang keras menjadi Key Map.

```go
// ❌ COMPILE ERROR: invalid map key type []string
// var invalidMap map[[]string]int 

// ✅ BENAR: Gunakan tipe primitive atau Struct/Array (comparable)
var validMap map[string]int
```

#### 🚫 Larangan 3: DILARANG Mengambil Alamat Elemen Map (*Unaddressable*)
Elemen map tidak punya alamat memori tetap karena bisa berpindah saat map bertumbuh (*rehashing*).

```go
m := map[string]int{"a": 10}

// ❌ COMPILE ERROR: cannot take the address of m["a"]
// ptr := &m["a"] 
```

#### 🚫 Larangan 4: DILARANG Memanggil Pointer Method pada Elemen Map
Karena unaddressable, pointer method tidak bisa dipanggil pada struct di dalam map.

```go
type User struct { Name string }
func (u *User) ChangeName(n string) { u.Name = n }

users := map[string]User{"admin": {Name: "Alice"}}
// ❌ COMPILE ERROR: cannot call pointer method ChangeName on User
// users["admin"].ChangeName("Bob")

// ✅ BENAR: Simpan POINTER di dalam map!
usersPtr := map[string]*User{"admin": &User{Name: "Alice"}}
usersPtr["admin"].ChangeName("Bob") // Sukses!
```

---

### 3. Structs

#### 💡 Karakteristik: Exported (Public) vs Unexported (Private)
* Huruf awal **KAPITAL** = Exported (Bisa diakses package luar).
* Huruf awal **kecil** = Unexported (Hanya bisa diakses di package lokal).

#### 🚫 Larangan: Huruf Kecil pada Field Struct untuk JSON Serialization
Paket `encoding/json` bekerja di luar package-mu. Field huruf kecil akan diabaikan!

```go
type Product struct {
    // ❌ SALAH: Tidak akan pernah di-export ke JSON (hasilnya kosong)!
    name  string `json:"name"` 

    // ✅ BENAR: Awali huruf kapital!
    Name  string `json:"name"` 
    Price int    `json:"price"`
}
```

---

## MODUL 3: Control Flow

### 💡 Karakteristik: Labeled Break & Go 1.22 Loop Scope

#### 🚫 Larangan 1: `break` Biasa di Dalam `switch` di Dalam `for`
`break` biasa hanya memutus `switch`, perulangan `for` akan tetap berlanjut!

```go
// ❌ SALAH: Loop tidak berhenti!
for {
    switch cmd {
    case "exit":
        break // Hanya keluar dari switch, loop tetap jalan terus!
    }
}

// ✅ BENAR: Gunakan Labeled Break
LoopUtama:
for {
    switch cmd {
    case "exit":
        break LoopUtama // Berhenti total dari for loop!
    }
}
```

#### 🚫 Larangan 2: DILARANG Memakai `defer` di Dalam Loop
`defer` baru dieksekusi saat **seluruh fungsi utama selesai**, bukan per iterasi loop!

```go
// ❌ BAHAYA: File descriptors bocor, error 'too many open files'!
for _, filename := range files {
    f, _ := os.Open(filename)
    defer f.Close() // Tertahan sampai fungsi selesai!
}

// ✅ BENAR: Bungkus dalam fungsi kecil agar defer langsung jalan di tiap iterasi
for _, filename := range files {
    func() {
        f, _ := os.Open(filename)
        defer f.Close() // Tutup file seketika fungsi anonim ini selesai!
    }()
}
```

---

## MODUL 4: Functions & Closures

### 💡 Karakteristik: 100% Pass-by-Value & Escape Analysis
Semua argumen ke fungsi selalu difotokopi. Jika pointer lokal di-return, compiler Go otomatis memindahkannya dari **Stack ke Heap**.

```go
// ✅ 100% Aman di Go (Escape Analysis otomatis memindahkan ke Heap)
func NewUser(name string) *User {
    return &User{Name: name} 
}
```

#### 🚫 Larangan 1: Mengabaikan Error
```go
// ❌ DILARANG di production:
data, _ := db.Query()

// ✅ BENAR: Selalu tangani error secara eksplisit
data, err := db.Query()
if err != nil {
    return fmt.Errorf("gagal query data: %w", err)
}
```

#### 🚫 Larangan 2: Hati-Hati Mutasi Nilai oleh Defer Closure
Fungsi `defer` yang berupa closure bisa menimpa **Named Return Value** sebelum nilai tersebut benar-benar dikembalikan ke pemanggil!

```go
func Calculate() (result int) {
    defer func() {
        result = 99 // Mengubah return value di menit-menit terakhir!
    }()
    return 10 // Yang sampai ke pemanggil adalah 99, bukan 10!
}
```

---

## MODUL 5: Pointers & Memory Management

### 💡 Karakteristik Utama: Operator `&` vs `*`
* `&x` (*Address-of*): Mengambil alamat memori.
* `*p` (*Dereference*): Datangi alamat dan akses/ubah isi rumahnya.

### 🚫 Larangan 1: DILARANG Dereference pada Pointer Nil
```go
var p *int // Nilai awal: nil

// ❌ RUNTIME PANIC: invalid memory address or nil pointer dereference
// *p = 50 

// ✅ BENAR: Alokasikan memori terlebih dahulu
val := 50
p = &val
*p = 50 // Aman!
```

### 🚫 Larangan 2: DILARANG Mengira `p = &val` Mengubah Data Asli Pemanggil
```go
// ❌ GAGAL: Hanya mengubah kertas alamat lokal 'p'
func Fail(p *int) {
    temp := 99
    p = &temp 
}

// ✅ BENAR: Gunakan tanda bintang untuk menimpa ISI RUMAHNYA!
func Success(p *int) {
    *p = 99 
}
```

### 🚫 Larangan 3: DILARANG Shallow Copy pada Struct dengan Pointer
```go
type Order struct {
    Customer *Customer // Field Pointer
}

order1 := Order{Customer: &Customer{Name: "Alice"}}

// ❌ BAHAYA (Shallow Copy): Kedua order menunjuk ke Customer yang sama!
order2 := order1
order2.Customer.Name = "Bob" 
fmt.Println(order1.Customer.Name) // Ikut berubah jadi "Bob"!

// ✅ BENAR (Deep Copy): Alokasikan rumah struct baru!
order2 := order1
order2.Customer = &Customer{
    Name: order1.Customer.Name,
}
order2.Customer.Name = "Bob" // order1 tetap aman bernilai "Alice"!
```

### 🚫 Larangan 4: "Pointer Abuse" (Semua Field Dijadikan Pointer)
Jadikan pointer **HANYA JIKA**:
1. Field opsional di JSON PATCH (`*int`, `*string`, `*bool`) agar bisa membedakan `nil` vs nilai `0`.
2. Kolom database yang bersifat `NULLABLE`.
3. Struct berukuran sangat besar demi efisiensi.

```go
// ✅ Trik Sakti Generic Pointer Helper:
func Ptr[T any](v T) *T { return &v }

req := UpdateProductRequest{
    Price: Ptr(100_000), // Langsung jadi *int tanpa variabel temp!
}
```

---

## MODUL 6: Methods & Interfaces

### 💡 Karakteristik: Duck Typing (Implicit Satisfaction)
Tidak butuh keyword `implements`. Cukup miliki method yang sama dengan interface.

```go
type Payer interface {
    Pay(amount int)
}

type Gopay struct{}
func (g Gopay) Pay(amount int) { fmt.Println("Bayar via Gopay:", amount) }

// Gopay otomatis dianggap sebagai Payer!
var p Payer = Gopay{} 
```

### 🚫 Larangan 1: DILARANG Value Receiver pada Struct dengan `sync.Mutex`
Value receiver akan memfotokopi gembok mutex. Gembok fotokopi yang terkunci tidak akan melindungi brankas asli $\rightarrow$ **DATA RACE!**

```go
type SafeWallet struct {
    mu      sync.Mutex
    Balance int
}

// ❌ SALAH: Passes lock by value!
func (w SafeWallet) DepositBad(amount int) {
    w.mu.Lock()
    defer w.mu.Unlock()
    w.Balance += amount
}

// ✅ BENAR: Selalu gunakan Pointer Receiver untuk Mutex!
func (w *SafeWallet) DepositGood(amount int) {
    w.mu.Lock()
    defer w.mu.Unlock()
    w.Balance += amount
}
```

### 🚫 Larangan 2: Aturan "Method Sets" pada Interface (`T` vs `*T`)
Jika method menggunakan Pointer Receiver `(t *T)`, kamu **DILARANG** mengoper Value biasa `T` ke Interface!

```go
type Worker interface { Work() }
type Person struct{}
func (p *Person) Work() {} // Pointer Receiver!

// ❌ COMPILE ERROR: Person does not implement Worker (Work method has pointer receiver)
// var w Worker = Person{} 

// ✅ BENAR: Wajib mengoper pointernya (&Person{})!
var w Worker = &Person{} 
```

### 🚫 Larangan 3: DILARANG Type Assertion Tanpa Comma-Ok
Type assertion **BUKAN** konversi tipe. Jika salah tipe, program akan langsung crash!

```go
var data any = "100" // String

// ❌ RUNTIME PANIC: interface conversion: interface {} is string, not int
// val := data.(int) 

// ✅ BENAR: Selalu gunakan Comma-Ok Idiom
val, ok := data.(int)
if !ok {
    fmt.Println("Bukan integer! Selamat dari panic!")
} else {
    fmt.Println("Hasil:", val)
}
```

### 🌟 Kaidah Emas Arsitektur Go:
> **"Accept Interfaces, Return Structs"**
> * **Parameter Input Menerima Interface:** Agar mudah diganti Mock DB saat Unit Test.
> * **Return Output Berupa Struct Konkret:** Agar pemanggil leluasa mengakses semua fiturnya tanpa terkunci interface.

```go
// ✅ Standar Industri:
func NewUserService(repo UserRepository) *UserService {
    return &UserService{repo: repo}
}
```
