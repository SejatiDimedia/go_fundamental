# Catatan Deep Dive: Introduction to Go & Internals

Dokumen ini merangkum seluruh hasil diskusi dan deep dive pada sub-topik **1. Language Basics > Introduction to Go**, mencakup arsitektur sistem, runtime internals, manajemen memori, hingga trade-off di production.

---

## 1. Latar Belakang & Filosofi Kelahiran Go

Go diciptakan pada tahun **2007** di Google oleh tiga tokoh legendaris: **Robert Griesemer**, **Rob Pike**, dan **Ken Thompson**.

### Masalah Nyata yang Ingin Diselesaikan di Google:
1. **Compilation Time yang Lambat (Dependency Hell di C++):**
   * Di C++, mekanisme `#include` menyebabkan preprosesor mem-parse ulang file header ribuan kali. Kompilasi sistem skala Google memakan waktu **45 menit hingga beberapa jam**.
   * Go memecahkannya dengan model kompilasi berbasis package terisolasi dan **melarang circular dependency** secara ketat (kompilasi hitungan milidetik/detik).
2. **Revolusi Hardware Multi-Core vs Heavy OS Threads:**
   * Hardware komputer modern bergeser ke CPU multi-core.
   * Concurrency tradisional di C++/Java bertumpu pada **OS Thread (Kernel Thread)** yang berat (alokasi memori 1–2 MB per thread, biaya *context switching* level kernel mahal).
   * Go merancang **Goroutine (Green Thread)** yang hanya butuh alokasi awal **~2 KB** di memori.
3. **Kompleksitas Sintaks yang Berlebihan:**
   * Fitur seperti inheritance bertingkat, template rumit, dan exception handling (`try-catch-finally`) membuat *control flow* sulit diprediksi saat dikerjakan oleh ribuan engineer.
   * Go mengambil pendekatan minimalis: hanya **25 keywords**, tidak ada inheritance (memakai *composition over inheritance*), dan error ditangani secara eksplisit sebagai *first-class value*.

---

## 2. Kompilasi & Single Binary: Apa Isi `go build`?

Berbeda dengan **C** (pure machine code tanpa GC/scheduler bawaan) dan **Java/Python** (membutuhkan JVM atau Interpreter terpisah), Go mengompilasi kode menjadi **Single Static Binary Executable**.

Di dalam satu file binary executable tersebut, sudah tersemat **Go Runtime**:
1. **Go Scheduler (GMP Model):** Mengatur eksekusi ratusan ribu goroutine di atas sedikit OS Thread.
2. **Garbage Collector (GC):** Pembersih memori otomatis bertipe *concurrent tri-color mark-and-sweep*.
3. **Netpoller:** Abstraksi asynchronous network I/O berbasis kernel OS (`epoll` di Linux, `kqueue` di macOS). Kode ditulis bergaya blocking/sinkron, tetapi di balik layar berjalan non-blocking.
4. **Memory Allocator:** Pengelola alokasi memori Stack dan Heap (berbasis TCMalloc).

### Trade-off di Dunia Production / Kubernetes:
* **Keuntungan:** Deployment sangat mudah (*zero external dependency*), bisa jalan di base image Docker minimalis (`scratch` / `alpine`), ukuran container image kecil (15–30 MB), dan startup time instan (tanpa warm-up JIT).
* **Kerugian:** Ukuran minimal binary Go lebih besar dibanding C (program Hello World Go berukuran ~1.5–2 MB vs C ~15 KB karena membawa Go Runtime). Di production backend besar, selisih ini tidak signifikan.

---

## 3. Concurrency Model: Heavy Threads vs Green Threads (Goroutines)

### Konsep Kernel Space vs User Space:
* **Kernel Space:** Wilayah inti OS yang mengontrol hardware secara langsung (CPU, RAM, Network Card). **OS Thread** hidup di sini. Context switch antar OS Thread mahal karena CPU harus bolak-balik antara User Mode dan Kernel Mode.
* **User Space:** Wilayah tempat aplikasi kita (Go app, browser, database) berjalan.

```text
┌────────────────────────────────────────────────────────┐
│ USER SPACE (Aplikasi Go)                               │
│                                                        │
│   Goroutine 1   Goroutine 2   Goroutine 3  ... (Ribuan)│
│        \             |             /                   │
│         ▼            ▼            ▼                    │
│        [ Go Runtime Scheduler (User-space Mandor) ]    │
└──────────────────────┬─────────────────────────────────┘
                       │ Memetakan M:N ke sedikit OS Thread
┌──────────────────────▼─────────────────────────────────┐
│ KERNEL SPACE (Sistem Operasi)                          │
│                                                        │
│        [ OS Thread 1 ]         [ OS Thread 2 ]         │
└──────────────────────┬─────────────────────────────────┘
                       ▼
                 [ CPU Cores ]
```

### Cara Kerja Go Runtime Scheduler (M:N Model):
* Kode pada akhirnya **tetap dieksekusi oleh OS Thread dan Core CPU fisik**.
* Bedanya: Go **membatasi jumlah OS Thread aktif** (default = jumlah core CPU, diatur via `GOMAXPROCS`).
* Di User Space, **Go Runtime Scheduler** bertindak seperti mandor: ia menggilir ratusan ribu Goroutine ke OS thread yang tersedia secara efisien tanpa merepotkan Kernel OS.

---

## 4. Manajemen Memori: Stack vs Heap & Escape Analysis

Di Go, memori RAM yang digunakan aplikasi terbagi menjadi dua wilayah utama:

### A. Stack Memory
* **Karakteristik:** Sangat cepat, berukuran kecil per-goroutine (mulai ~2 KB), bersifat LIFO (Last In First Out).
* **Siklus Hidup:** Variabel lokal di stack **langsung musnah seketika** saat fungsi selesai dieksekusi. Bebas overhead GC.

### B. Heap Memory
* **Karakteristik:** Area penyimpanan bersama yang ukurannya dinamis dan lebih lambat dibanding Stack.
* **Siklus Hidup:** Data di Heap tetap hidup meskipun fungsi pembuatnya sudah selesai. Data ini **wajib dibersihkan oleh Garbage Collector** jika sudah tidak ada pointer yang merujuk kepadanya.

### C. Miskonsepsi Pointer: "Apakah Semua Pointer Masuk ke Heap?"
**JAWABANNYA: TIDAK.**  
Yang menentukan alokasi ke Heap atau Stack adalah **Escape Analysis** oleh compiler Go.

1. **Pointer TIDAK Kabur (Tetap di STACK):**
   ```go
   func proses() {
       x := 42
       update(&x) // Pointer dipinjamkan ke fungsi anak
       println(x)
   }
   func update(p *int) {
       *p = 100 // update() selesai duluan sebelum proses() selesai. Stack 'proses' masih hidup.
   }
   ```
   *Analisis Compiler:* Data `x` **tidak kabur** melampaui stack pembuatnya. Memori `x` tetap di **Stack**.

2. **Pointer Kabur (Escape to HEAP):**
   ```go
   func buatData() *int {
       y := 100
       return &y // Alamat memori &y di-return ke atas/ke pemanggil
   }
   ```
   *Analisis Compiler:* Jika `y` ditaruh di stack, saat `buatData()` selesai, memori `y` akan musnah dan pemanggil memegang pointer liar (*dangling pointer*). Compiler secara otomatis **memindahkan `y` ke Heap**.

> **Tips Senior:** Jalankan perintah `go build -gcflags="-m"` di terminal untuk melihat laporan Escape Analysis compiler Go secara langsung.

---

## 5. Garbage Collector (GC) & Scavenging

1. **Kapan GC Bekerja?**
   GC berjalan secara berkala dan konkuren di background **selama aplikasi sedang aktif** untuk mencari dan membebaskan data Heap yang sudah tidak bisa dijangkau (*unreachable objects*).
2. **Tugas-tugas GC:**
   * **Tracing:** Menandai objek hidup vs mati dengan teknik *tri-color marking*.
   * **Sweeping & Recycling:** Mengumpulkan memori bebas ke allocator internal Go agar bisa digunakan kembali tanpa minta RAM baru ke OS.
   * **Scavenging:** Mengembalikan memori yang lama menganggur ke OS Kernel (via syscall `madvise(MADV_DONTNEED)` di Linux) agar RAM server tidak mubazir.
3. **Environment Variables Penting:**
   * `GOMEMLIMIT`: Batasan lunak penggunaan RAM sebelum GC bekerja lebih agresif. Mencegah container terkena *OOMKilled* di Kubernetes.
   * `GOGC`: Target persentase pertumbuhan heap untuk memicu siklus GC (default: `100`).
   * `GOMAXPROCS`: Jumlah maksimal OS Thread yang mengeksekusi kode Go secara paralel (default: jumlah logical CPU cores).

---

## 6. Circular Dependency (Import Cycle)

### Mengapa Compiler Go Melarang Circular Import?
Jika `package A` meng-import `package B`, dan `package B` juga meng-import `package A`:
* Compiler tidak bisa membentuk grafik dependensi terarah (DAG: *Directed Acyclic Graph*).
* Menandakan desain kode buruk (*tight coupling* / spaghetti architecture).

### Cara Mengatasinya:
1. **Shared / Domain Package:** Pindahkan entitas/struct bersama ke package yang lebih rendah (misal `model` atau `entity`).
2. **Inversion of Control (Interfaces):** Definisikan interface di consumer package.
3. **Layered / Clean Architecture:** Dependensi satu arah:
   `Handler` ──> `Service / Usecase` ──> `Repository` ──> `Domain / Model`.

---

## 7. Entry Point & Compiler Strictness

1. **`package main` dan `func main()`:**
   * Program executable wajib memiliki `package main` dan `func main()` (semua huruf kecil, tanpa parameter, tanpa return value).
   * Menamai fungsi menjadi `func Main()` dengan huruf kapital akan menyebabkan compiler/linker error:  
     `function main is undeclared in the main package`.
2. **Disiplin Compiler Go:**
   * Go menolak kompilasi jika ada variabel atau import yang dideklarasikan tetapi tidak dipakai (*unused variable / unused import*).
   * Filosofi: Mencegah *dead code*, menjamin kecepatan build, dan mempermudah keterbacaan kode dalam tim besar.
