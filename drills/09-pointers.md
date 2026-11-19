# Catatan Deep Dive: Pointers & Memory Management di Go

Dokumen ini merangkum seluruh hasil diskusi dan deep dive pada sub-topik **5. Pointers**, mencakup operator `&` dan `*`, nil pointer dereference, pelarangan pointer arithmetic, constructor pattern, serta trade-off performa Stack vs Heap.

---

## 1. Operator Dasar: `&` (Address-of) vs `*` (Dereference)

* **`&x` (Address-of):** Mengambil alamat memori tempat variabel `x` disimpan (menghasilkan pointer `*T`).
* **`*p` (Dereference):** Mendatangi alamat memori yang disimpan di `p` untuk membaca atau mengubah nilai aslinya.

```go
x := 42
p := &x  // p menyimpan alamat x
*p = 99  // x sekarang bernilai 99
```

---

## 2. The Nil Pointer Dereference Trap

Zero value dari setiap tipe pointer di Go adalah **`nil`**.
* Mencoba membaca atau mengisi field dari pointer yang bernilai `nil` akan memicu:
  🔥 **`panic: runtime error: invalid memory address or nil pointer dereference`**
* **Solusi:** Selalu alokasikan memori struct terlebih dahulu (menggunakan `&MyStruct{}` atau `new(MyStruct)`) sebelum mengakses field-nya.

---

## 3. Mengapa Go Melarang Pointer Arithmetic (`p++`)?

Di C/C++, programmer bisa menggeser pointer sembarangan (`p++`).
* **Go melarangnya demi Memory Safety:** Mencegah kebocoran data rahasia di RAM, korupsi memori, dan eksploitasi keamanan (*buffer overflow*).
* Untuk menelusuri memori berurutan, Go menyediakan **Slice** yang aman dan dilengkapi pemeriksaan batas (*bounds checking*).

---

## 4. Constructor Pattern & Escape Analysis

```go
func NewAccount(id string, balance int) *Account {
    return &Account{ID: id, Balance: balance}
}
```

* **Di C/C++:** Mengembalikan pointer ke variabel lokal adalah bug fatal (*Dangling Pointer*) karena memori Stack fungsi langsung hancur.
* **Di Go (100% Aman):** Compiler Go menjalankan **Escape Analysis**. Begitu melihat alamat struct di-return ke luar, Go otomatis memindahkan alokasinya ke **Heap Memory**, sehingga aman dari kehancuran stack.

---

## 5. Kapan Pointer Justru Memperlambat Performa?

| Skenario | Lebih Cepat Pakai | Alasan Teknis |
| :--- | :--- | :--- |
| **Struct Kecil ($\le$ 64 byte)** | **Value** | Disimpan di **Stack** (0 beban Garbage Collector) dan memanfaatkan kecepatan CPU L1/L2 Cache (tanpa *pointer chasing*). |
| **Struct Besar (banyak field / buffer)** | **Pointer** | Menghindari biaya meng-copy puluhan/ratusan byte memori berulang kali. |
| **Butuh Mutasi Data Asli** | **Pointer** | Satu-satunya cara mengubah data pemanggil (Go is 100% pass-by-value). |
| **Read-Only Data** | **Value** | Menjamin data tidak termutasi dan ramah GC. |
