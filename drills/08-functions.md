# Catatan Deep Dive: Functions di Go

Dokumen ini merangkum seluruh hasil diskusi dan deep dive pada sub-topik **4. Functions**, mencakup filosofi penolakan exception, named return values, fungsi variadic, mitos pass-by-reference, serta closure & escape analysis.

---

## 1. Multiple Return Values vs Exceptions (`try-catch`)

Go menolak mekanisme `try-catch` / exceptions karena:
1. **Errors are Values:** Error dalam backend (jaringan, database, disk) adalah kejadian lumrah (*ordinary control flow*), bukan bencana luar biasa (*exceptional*).
2. **Eksplisit (No Invisible GOTO):** Exception di Java/Python melompat tak terlihat melewati banyak stack frame. Di Go, `val, err := f()` membuat kegagalan terlihat jelas dan wajib ditangani di tempat.
3. **Performa:** Mengembalikan nilai error di Go nyaris berbiaya nol (menggunakan register CPU), dibanding membuat stack trace exception yang mahal.

---

## 2. Named Return Values & Naked Returns

```go
func Calculate() (result int, err error) {
    result = 10
    return // Naked Return: otomatis return result, err
}
```

* **Manfaat:** Berfungsi sebagai dokumentasi diri (*self-documenting code*) di godoc.
* **Perilaku Return + Defer:**
  1. `return 5` meng-assign nilai ke variabel named return (`result = 5`).
  2. Fungsi `defer` dijalankan dan **bisa memodifikasi named return tersebut** sebelum benar-benar keluar ke pemanggil!
* **Best Practice:** Hindari *naked return* pada fungsi yang panjang (> 10 baris) karena membingungkan pembaca dan rawan bug *variable shadowing*. Selalu gunakan return eksplisit: `return result, err`.

---

## 3. Variadic Functions (`...T`) & Slice Unpacking (`s...`)

1. Parameter `...T` diperlakukan sebagai slice `[]T` di dalam fungsi.
2. Memanggil dengan elemen lepas: `Sum(1, 2, 3)`.
3. Memanggil dengan slice yang sudah ada: **Wajib dibongkar** dengan ellipsis: `Sum(mySlice...)`.
4. **Efek Memori:** Mengoper `mySlice...` ke parameter variadic **membagikan Underlying Array yang sama**. Modifikasi elemen di dalam fungsi akan merubah slice asli pemanggil!

---

## 4. Mitos "Pass by Reference" di Go

> **FAKTA MUTLAK: Go adalah 100% Pass-by-Value!**

Tidak ada satu pun tipe data di Go yang dioper secara *pass-by-reference*:
* **Slice:** Yang di-copy adalah *Slice Header* (24 byte).
* **Map:** Yang di-copy adalah *Pointer Header* (`*hmap`).
* **Pointer (`*T`):** Yang di-copy adalah *Alamat Memori* (8 byte).
* **Struct / Int / Array:** Yang di-copy adalah seluruh nilainya (*pure value copy*).

---

## 5. Closures & Escape Analysis (Stack vs Heap)

Closure adalah fungsi anonim yang mengunci (*close over*) variabel di sekelilingnya.

```go
func CreateCounter() func() int {
    count := 0
    return func() int {
        count++
        return count
    }
}
```

* **Kenapa `count` tidak mati saat `CreateCounter` selesai?**
  Compiler Go menjalankan **Escape Analysis**: karena `count` masih dirujuk oleh fungsi anonim di luar, compiler memindahkan alokasi `count` dari **Stack** ke **Heap** (*count escapes to heap*).
* Variabel di Heap tetap hidup selama closure masih memegangnya, dan baru dibersihkan oleh Garbage Collector saat sudah tidak ada referensi.

---

## 6. Timing Evaluasi Argumen pada `defer`

1. **Argumen Langsung (`defer fmt.Println(a)`):**
   * Nilai argumen dievaluasi dan dikunci **pada saat baris `defer` dipanggil**.
   * Jika `a` diubah setelahnya, nilai yang dicetak tetap nilai saat baris `defer` dibaca.
2. **Closure (`defer func() { fmt.Println(b) }()`):**
   * Tidak ada argumen yang dikunci. Variabel `b` ditangkap oleh referensi.
   * Nilai `b` baru dibaca **saat fungsi anonim dieksekusi di akhir fungsi**.

---

## 7. Functions as First-Class Citizens (Callbacks)

Fungsi di Go dapat dioper sebagai parameter:
* `myFunc` (tanpa tanda kurung): Mengoper referensi objek fungsinya.
* `myFunc()` (dengan tanda kurung): Menjalankan fungsinya sekarang juga dan menghasilkan return value-nya.
* Saat mengoper callback ke parameter bertipe `func(...)`: **Hanya kirimkan nama fungsinya tanpa tanda kurung**.
