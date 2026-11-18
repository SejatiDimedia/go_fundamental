# Catatan Deep Dive: Arrays & Slices

Dokumen ini merangkum seluruh hasil diskusi dan deep dive pada sub-topik **2. Composite Types > Arrays & Slices**, mencakup perbedaan array vs slice, alokasi memori `make`, pertumbuhan kapasitas, hingga memory leak pada slice.

---

## 1. Array vs Slice: Perbedaan Fundamental

| Karakteristik | Array (`[N]T`) | Slice (`[]T`) |
| :--- | :--- | :--- |
| **Ukuran** | Fixed (Ukuran adalah bagian dari tipe data, misal `[3]int != [5]int`). | Dinamis (Panjang bisa bertambah via `append`). |
| **Struktur Memori** | Menyimpan seluruh elemen langsung di tempat alokasi. | Struct 24-byte (Pointer ke Underlying Array, Len, Cap). |
| **Passing ke Fungsi** | **Pure Value Copy:** Seluruh elemen diduplikasi ke stack frame baru. Mutasi di dalam fungsi tidak mempengaruhi data luar. | **Header Copy:** Yang dicopy hanya header 24-byte. Mutasi elemen via indeks (`s[0] = 999`) **mengubah Underlying Array asli di luar**. |

---

## 2. Anatomi `make([]T, len, cap)` & Jebakan Index

```go
s := make([]int, 3, 5)
```
* **`Len = 3`**: Ada 3 elemen aktif yang langsung diisi zero value: `[0, 0, 0]`.
* **`Cap = 5`**: Underlying array menyiapkan 5 slot memori (ada 2 slot cadangan kosong).
* **Awas Panic:** Kamu hanya bisa mengakses langsung `s[0]`, `s[1]`, `s[2]`. Mengakses `s[3]` akan memicu **panic: index out of range** meskipun `cap = 5`.
* **`make([]int, 0, 5)`:** Sangat dianjurkan untuk penampung data kosong yang akan di-`append`. Akses `s[0]` langsung akan panic jika belum di-`append`.

---

## 3. Slice Growth (Pertumbuhan Kapasitas di Balik Layar)

Ketika `append()` dipanggil pada slice yang kapasitasnya sudah penuh (`Len == Cap`):
1. Go Runtime (`runtime.growslice`) mengalokasikan **Underlying Array BARU** yang lebih besar di RAM.
2. Seluruh data dari array lama disalin ke array baru.
3. Elemen baru ditambahkan ke array baru.
4. Pointer slice diarahkan ke array baru.
5. **Dampaknya:** Slice yang telah tumbuh kini terlepas (*decoupled*) dari Underlying Array lama.

---

## 4. Rumus Slicing `[low:high]`

Slicing bersifat *half-open interval*:
$$\text{Panjang Element} = \text{high} - \text{low}$$
```go
base := []int{10, 20, 30, 40, 50}
sub := base[1:3] // Mengambil index 1 dan 2 -> [20, 30]
```
* `sub[0]` menunjuk ke `base[1]`.
* Mengubah `sub[0] = 999` otomatis mengubah `base[1] = 999` karena keduanya berbagi Underlying Array yang sama.

---

## 5. Memory Leak Trap pada Slice (Kasus File/Payload Besar)

### Masalah:
Memotong slice kecil dari slice raksasa:
```go
hugeLog := readHugeFile() // 1 GB di RAM
status := hugeLog[:10]   // hanya butuh 10 byte
return status
```
Karena `status` masih memegang pointer ke Underlying Array 1 GB tersebut, **Garbage Collector TIDAK BISA membersihkan 1 GB RAM tersebut**.

### Solusi Standar Industri:
Putuskan hubungan dengan membuat alokasi baru yang terisolasi:

```go
// Cara 1: Gunakan copy()
cleanStatus := make([]byte, 10)
copy(cleanStatus, hugeLog[:10])
return cleanStatus

// Cara 2: Gunakan slices.Clone() (Go 1.21+)
import "slices"
return slices.Clone(hugeLog[:10])
```
*Dengan cara ini, memori 1 GB lama bisa langsung dibersihkan oleh Garbage Collector.*


---

## 6. Full Slice Expression (3-Index Slicing: `s[low:high:max]`)

### Masalah pada Slicing 2-Indeks Biasa:
```go
base := []int{10, 20, 30, 40, 50}
sub := base[1:3] // len: 2, cap: 4 (mewarisi sisa kapasitas sampai ujung)

sub = append(sub, 999) 
// KARENA KAPASITAS MASIH CUKUP, append TIDAK MEMBUAT ARRAY BARU!
// Elemen base[3] yang tadinya 40 TERTIMPA menjadi 999!
```

### Solusi: Kunci Kapasitas dengan Indeks ke-3 (`:max`)
$$\text{Len} = \text{high} - \text{low}$$
$$\text{Cap} = \text{max} - \text{low}$$

```go
sub := base[1:3:3] // len: 2, cap: 2
sub = append(sub, 999) 
// Karena len == cap, Go DIPAKSA mengalokasikan Underlying Array baru!
// Variabel base 100% aman dari kerusakan data tertimpa.
```

---

## 7. Mengapa `append` di Dalam Fungsi Tidak Mengubah `Len` di Pemanggil?

Ketika slice dioper ke fungsi, **Slice Header** (struct 24 byte: `{Data, Len, Cap}`) di-copy secara *pass by value*.
* Jika kapasitas masih cukup, `append` di dalam fungsi **BENAR-BENAR MENULIS** elemen baru ke Underlying Array bersama.
* Namun, perubahan nilai `Len` hanya terjadi pada Slice Header lokal di dalam fungsi tersebut.
* Di fungsi pemanggil (`main`), nilai `Len` pada header pemanggil **tidak berubah**. Pemanggil hanya membaca sampai batas `Len` lamanya.
* **Solusi Idiomatic:** Kembalikan slice hasil append (`return append(...)`) atau kirim pointer ke slice (`*[]T`).


---

## 8. Teknik Hapus Elemen Slice di Industri: Ordered vs Fast Delete

### A. Ordered Delete ($O(N)$ - Menjaga Urutan)
```go
func RemoveOrdered(s []int, index int) []int {
    return append(s[:index], s[index+1:]...)
}
```
* **Kapan Pakai:** Wajib dipakai ketika urutan data penting (antrean FIFO, urutan transaksi timestamp, leaderboard).
* **Kekurangan:** Lambat untuk data besar karena harus menggeser seluruh elemen di belakangnya ke kiri.

### B. Fast Unordered Delete ($O(1)$ - Kecepatan Maksimal)
```go
func RemoveFast(s []int, index int) []int {
    s[index] = s[len(s)-1] // Timpa dengan elemen paling belakang
    return s[:len(s)-1]   // Potong ekor
}
```
* **Kapan Pakai:** Dipakai ketika urutan data tidak penting (daftar worker pool, connection set).
* **Kelebihan:** Eksekusi instan $O(1)$ berapapun jumlah datanya di RAM.
