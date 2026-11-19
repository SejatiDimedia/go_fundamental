# Catatan Deep Dive: Control Flow di Go

Dokumen ini merangkum seluruh hasil diskusi dan deep dive pada sub-topik **3. Control Flow**, mencakup `if` short statement, `switch` tagless & fallthrough, labeled break, evolusi loop variable Go 1.22+, serta iterasi string byte vs rune.

---

## 1. `if` dengan Short Statement & Scope Isolation

```go
if err := processData(); err != nil {
    return err
}
// Di sini 'err' sudah mati (out of scope)
```

* **Manfaat:** Mencegah *variable pollution* (pencemaran variabel) dan *shadowing bugs*. Variabel yang diinisialisasi pada short statement hanya hidup di dalam blok `if` dan `else`.

---

## 2. Fitur-Fitur Unik `switch` di Go

1. **Otomatis `break`:** Tidak perlu menulis keyword `break` di akhir setiap case. Go otomatis keluar dari switch begitu ada case yang cocok.
2. **Keyword `fallthrough`:** Memaksa eksekusi lanjut ke blok case di bawahnya tanpa mengecek kondisi case tersebut.
3. **Tagless Switch (Pengganti `if - else if`):**
   ```go
   switch {
   case score >= 90:
       grade = "A"
   case score >= 80:
       grade = "B"
   default:
       grade = "C"
   }
   ```
4. **Multiple Expression dalam 1 Case:**
   ```go
   case "A", "B", "C":
       fmt.Println("Lulus")
   ```

---

## 3. Labeled Break (Jebakan `for` + `switch`)

Perintah `break` polos di dalam `case` sebuah `switch` **HANYA menghentikan `switch`**, BUKAN menghentikan `for` loop pembungkusnya!

### Solusi: Beri Label pada Loop
```go
WorkerLoop:
for {
    switch status {
    case "STOP":
        break WorkerLoop // ✅ Keluar dari for loop
    }
}
```

---

## 4. Evolusi Semantik Loop Variable (Go 1.22+)

### Sebelum Go 1.22 (Go 1.21 ke bawah):
* Variabel `v` pada `for _, v := range slice` **hanya dialokasikan 1 kali di 1 alamat memori**.
* Goroutine atau closure di dalam loop sering mencetak nilai terakhir yang sama (`30 30 30`).

### Mulai Go 1.22+:
* Setiap putaran iterasi **mengalokasikan variabel `v` baru** di alamat memori yang berbeda.
* Closure dan goroutine di dalam loop sekarang aman dan mencetak nilai yang diharapkan (`10 20 30`).

---

## 5. Iterasi String: Byte vs Rune (UTF-8 Traps)

1. `len(str)` mengembalikan **jumlah byte**, bukan jumlah karakter.
2. **Loop Byte (`for i := 0; i < len(s); i++`)**: Mengakses `s[i]` per-byte (`uint8`). Karakter multi-byte seperti emoji atau huruf non-Latin (Arab, Jepang, Mandarin) akan terpotong dan menghasilkan karakter rusak (*mojibake*).
3. **Loop Rune (`for idx, r := range s`)**: Go mendekode byte-byte menjadi `rune` (Unicode code point). Karakter emoji tetap utuh, dan `idx` melompat sesuai jumlah byte yang dikonsumsi oleh rune tersebut.

---

## 6. The "Defer in Loop" Pitfall (File Descriptor Leak)

Keyword `defer` terikat pada **Function Scope**, bukan Block/Loop Scope.
* Menulis `defer file.Close()` di dalam loop 10.000 file akan membuat seluruh file tetap terbuka sampai fungsi selesai.
* Memicu OS error: `too many open files` (*File Descriptor Exhaustion*).
* **Solusi:** Bungkus isi loop ke dalam fungsi anonim `func() { ... defer file.Close() }()` agar `defer` dipanggil di setiap akhir iterasi.

---

## 7. Evaluasi Snapshot pada `for range` Slice

Ekspresi `range slice` dievaluasi **satu kali saja di awal loop**.
* Batas iterasi dikunci ke `len(slice)` awal.
* Jika slice di-`append` di dalam loop, elemen baru yang ditambahkan **tidak akan di-loop** pada perulangan yang sedang berjalan. Mencegah infinite loop tak disengaja.

