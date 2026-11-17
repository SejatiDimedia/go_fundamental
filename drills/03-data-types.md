# Catatan Deep Dive: Data Types, Strings, & UTF-8

Dokumen ini merangkum seluruh hasil diskusi dan deep dive pada sub-topik **1. Language Basics > Data Types**, mencakup strict typing, anatomi string, rune, UTF-8, hingga representasi nilai uang.

---

## 1. Strict Typing & Type Conversion

Go adalah bahasa yang bertipe statis kuat (*strongly typed*).
* **No Implicit Type Conversion:** Go sama sekali tidak melakukan konversi tipe data otomatis (auto-casting).
* Menjumlahkan `int` dan `float64` atau `int` dan `int64` akan memicu **Compile Error: mismatched types**.
* **Wajib Eksplisit:** Programmer wajib menulis konversi secara sadar: `int64(a) + b`.
* **Bahaya Downcasting (Overflow):**  
  Mengonversi integer berukuran besar ke kecil (misal `int16(300)` ke `int8`) tidak memicu runtime error, melainkan **pemotongan bit (*truncation*)** secara diam-diam (*silent overflow*):
  `300` (16-bit) -> `44` (8-bit).

---

## 2. Anatomi Teks di Go: `byte` vs `rune` vs `string`

### A. `byte`
* Alias resmi untuk **`uint8`** (angka 0–255).
* Berukuran **1 Byte (8 bit)**.
* Standar untuk karakter ASCII Latin dasar ('A'-'Z', '0'-'9').

### B. `rune`
* Alias resmi untuk **`int32`** (angka integer 4 byte).
* Merepresentasikan **1 Unicode Code Point** (karakter tunggal di dunia: huruf Kanji, Arab, Emoji, aksara daerah).
* Menggunakan tanda kutip tunggal: `'A'`, `'™'`, `'🚀'`.

### C. `string`
* Di Go, string secara internal adalah **read-only slice of bytes (`[]byte`)**.
* **Immutable:** Karakter di dalam string tidak bisa diubah langsung (`str[0] = 'X'` dilarang compiler).
* Menggunakan tanda kutip ganda: `"Hello"`.

---

## 3. Panjang Byte vs Panjang Karakter (UTF-8)

Standar UTF-8 bersifat *variable-length* (panjang byte per karakter bervariasi antara 1 hingga 4 byte):
* Huruf `'G'` = 1 byte
* Huruf `'o'` = 1 byte
* Simbol `'™'` = 3 bytes

```go
str := "Go™"
len(str)                        // 5 (Menghitung jumlah BYTE)
utf8.RuneCountInString(str)     // 3 (Menghitung jumlah KARAKTER / Runes)
```

### Cara Manipulasi String yang Benar:
Karena string bersifat *immutable*, manipulasi string dilakukan dengan konversi ke slice of rune:
```go
runes := []rune(str)
runes[0] = 'J'
result := string(runes) // Membuat string baru hasil copy
```

---

## 4. Larangan Menggunakan `float` untuk Nilai Uang (Financial / Payment)

### Masalah: Standar Biner IEEE 754
Pecahan desimal seperti `0.1` atau `0.2` tidak bisa direpresentasikan secara tepat dalam biner komputer.  
Hasil dari `0.1 + 0.2` di Go adalah:
`0.30000000000000004`

### Solusi Standar Industri:
1. **Gunakan `int64` dengan Satuan Terkecil (*Smallest Currency Unit*):**
   * Rupiah: Simpan bulat dalam Rupiah (`Rp 50.000` -> `50000`).
   * USD: Simpan dalam Cents (`$10.50` -> `1050`).
   * Crypto: Simpan dalam Satoshi / Wei.
   * *Nol risiko pembulatan, operasi cepat dan eksak.*
2. **Gunakan Library `shopspring/decimal`:**
   * Untuk kalkulasi suku bunga bank atau pembagian persen rumit yang wajib presisi di level desimal.
