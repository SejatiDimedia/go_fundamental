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


---

## 6. Mengubah Isi Rumah (`*p = val`) vs Mengubah Kertas Alamat (`p = &val`)

* **`*p = val` (Dereference):** Mengunjungi alamat memori yang ditunjuk dan menimpa nilai aslinya. Variabel asli di pemanggil ikut berubah.
* **`p = &val` (Reassignment):** Hanya mengubah catatan pointer lokal milik fungsi itu sendiri. Variabel asli di pemanggil **TIDAK** terpengaruh sama sekali.

---

## 7. Slice Header & Kapan Butuh Pointer to Slice (`*[]T`)

Slice di Go adalah struct 24-byte (Slice Header) yang berisi:
1. `Data unsafe.Pointer` (pointer ke underlying array di RAM)
2. `Len int` (panjang aktif)
3. `Cap int` (kapasitas maksimum)

* **`modifyElement(s []int)`:** Mengubah `s[0] = 88` berhasil tembus ke pemanggil karena pointer `s.Data` menunjuk ke array yang sama.
* **`appendElement(s []int)`:** GAGAL menambah panjang slice di pemanggil karena Slice Header di-copy by value (`s.Len` pemanggil tidak terupdate, dan jika kapasitas penuh, alokasi array baru hanya disimpan di kartu lokal `s`).
* **Solusi Mutasi Slice:** Gunakan pointer to slice `func Pop(s *[]int)` dan dereference dengan tanda kurung `*s = (*s)[:lastIndex]`.

---

## 8. Shallow Copy vs Deep Copy pada Struct

Ketika sebuah struct berisi field bertipe pointer (misal `Customer *Customer`):
* **Shallow Copy (`order2 := order1`):** Hanya menyalin nilai field langsung. Pointer `order2.Customer` akan mencatat **alamat RAM yang sama** dengan `order1.Customer`. Mengubah `order2.Customer.Name` akan **merusak** `order1.Customer.Name`!
* **Deep Copy (Solusi):** Alokasikan struct baru secara eksplisit di memori terpisah:
  ```go
  order2 := order1
  order2.Customer = &Customer{
      Name: order1.Customer.Name,
  }
  ```


---

## 9. Pointer to Pointer (`**T`)

* Variabel yang menyimpan alamat memori dari variabel pointer lain.
* Digunakan ketika kita ingin sebuah fungsi mengalihkan ke mana pointer pemanggil menunjuk:
  ```go
  func allocate(pp **int) {
      val := 42
      *pp = &val // Mengisi variabel pointer pemanggil dengan alamat val
  }
  ```

---

## 10. Pointer untuk Field Opsional (JSON PATCH & Database Pattern)

* Zero value dari tipe data primitif (`int` = 0, `string` = `""`, `bool` = `false`).
* Zero value dari SEMUA tipe pointer adalah **`nil`**.
* **Masalah:** Jika menggunakan `Price int`, backend tidak bisa membedakan antara *"klien tidak mengirim harga"* vs *"klien mengirim harga 0 (promo gratis)"*.
* **Solusi Industri:** Gunakan pointer `Price *int`:
  * `req.Price == nil` $\rightarrow$ Klien tidak mengirim field $\rightarrow$ Jangan ubah data di database.
  * `req.Price != nil` $\rightarrow$ Klien sengaja mengirim nilai (meskipun 0) $\rightarrow$ Update database dengan `*req.Price`.
