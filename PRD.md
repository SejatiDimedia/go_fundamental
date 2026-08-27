# PRD Lengkap: Kurikulum Pembelajaran Golang — Fundamental sampai Mahir (Job-Ready)

**Versi:** 2.0 (Lengkap)
**Tanggal:** 26 Agustus 2026
**Target akhir:** Siap bekerja sebagai Golang Developer level Middle-Senior, fokus e-commerce/payment system
**Estimasi durasi:** 16 minggu (± 2-3 jam/hari, fleksibel disesuaikan kecepatan belajar)

---

## 1. Ringkasan Eksekutif

Dokumen ini adalah kurikulum menyeluruh untuk mempelajari Golang dari nol/penyegaran hingga level siap kerja di posisi middle-senior, dengan spesialisasi ke arah sistem e-commerce dan payment. Kurikulum dipecah menjadi 8 fase, masing-masing punya materi rinci, output nyata (project/kode), dan kriteria selesai (Definition of Done).

## 2. Tujuan Program

| # | Tujuan | Indikator Sukses |
|---|--------|-------------------|
| 1 | Menguasai seluruh syntax dasar-lanjutan Go | Bisa menulis kode tanpa membuka referensi untuk kasus umum |
| 2 | Memahami konkurensi Go secara mendalam | Bisa membuat sistem concurrent yang aman dari race condition |
| 3 | Mampu merancang & membangun REST API skala production | Punya API dengan auth, validasi, testing, dokumentasi |
| 4 | Menguasai database SQL & NoSQL dari sisi backend Go | Paham transaction, indexing, caching, migration |
| 5 | Memahami arsitektur microservices | Bisa memecah monolith jadi beberapa service yang saling komunikasi |
| 6 | Menguasai pola-pola sistem pembayaran | Paham idempotency, distributed transaction, reconciliation |
| 7 | Siap deployment & observability | Bisa containerize, monitoring dasar, CI/CD sederhana |
| 8 | Siap interview teknis level middle-senior | Lulus simulasi live coding & system design |

## 3. Prasyarat

- Familiar dasar pemrograman (variabel, fungsi, logika)
- Komputer dengan Go terinstall (versi terbaru stabil)
- Editor (VS Code/GoLand) + akses terminal
- Akses internet untuk dokumentasi resmi (go.dev) dan referensi

## 4. Struktur Kurikulum (8 Fase)

---

### FASE 1: Fundamental Syntax (Minggu 1-2)

**Tujuan:** Membangun pondasi sintaks Go yang kokoh.

**Materi rinci:**
- Instalasi Go, struktur project, `go run` vs `go build`
- Package & import, `package main`, `func main()`
- Variabel: `var`, `:=`, zero value, scope (block, package, global)
- Tipe data dasar: `int`, `int8/16/32/64`, `uint`, `float32/64`, `string`, `bool`, `byte`, `rune`
- Konstanta (`const`) dan `iota` untuk enumerasi
- Type conversion & type inference
- Operator: aritmatika, logika, perbandingan, bitwise
- Percabangan: `if/else`, `if` dengan inisialisasi, `switch` (termasuk switch tanpa kondisi, fallthrough)
- Perulangan: `for` klasik, `for` sebagai while, infinite loop, `range`
- `break`, `continue`, label pada loop
- Array vs Slice (perbedaan fundamental: fixed-size vs dynamic, cara kerja underlying array)
- Operasi slice: `append`, `copy`, slicing `[low:high]`, `make`
- Map: deklarasi, CRUD, cek keberadaan key (`value, ok := m[key]`), iterasi map (urutan tidak terjamin)
- Function: parameter, multiple return value, named return, variadic function
- Function sebagai first-class citizen (function type, anonymous function, closure)
- `defer` — cara kerja dan use case umum (cleanup resource)
- String manipulation dasar (`strings`, `strconv` package)
- Formatting output (`fmt.Println`, `Printf`, `Sprintf`, verb `%v %d %s %T` dll)

**Deliverable:**
- 5 program CLI kecil: kalkulator, pengelola daftar belanja (pakai slice & map), konversi suhu, FizzBuzz varian, pengecek palindrome
- Semua program disimpan di repo GitHub `golang-fundamentals`

**Definition of Done:**
- Bisa menulis ulang 80% program tanpa membuka dokumentasi
- Paham betul perbedaan array-slice dan value-reference semantics

---

### FASE 2: Struct, Interface, Error Handling, Pointer (Minggu 3-4)

**Tujuan:** Memahami konsep OOP ala Go dan filosofi error handling.

**Materi rinci:**
- Struct: deklarasi, nested struct, anonymous struct
- Method: receiver value vs pointer receiver (kapan pakai yang mana)
- Struct tag (untuk JSON, validasi, dsb)
- Embedding struct (composition over inheritance)
- Interface: deklarasi, implicit implementation, empty interface (`any`)
- Type assertion & type switch
- Pointer: `&`, `*`, kapan perlu pointer, pointer ke struct
- Error handling idiomatic: `if err != nil`
- Membuat custom error: `errors.New`, `fmt.Errorf` dengan `%w` (error wrapping)
- `errors.Is` dan `errors.As` untuk cek tipe error
- Panic, recover, dan kapan (jarang) digunakan
- Package `sort` untuk custom sorting dengan interface `sort.Interface`
- Generics dasar (Go 1.18+): type parameter, constraint, generic function/struct

**Deliverable:**
- Project "Manajemen Inventaris Toko" (CLI): struct Produk, interface Pembayaran (Cash/Transfer), custom error untuk stok habis
- Refactor project Fase 1 dengan struct & interface

**Definition of Done:**
- Kode terorganisir dalam beberapa file/package
- Tidak ada unhandled error atau panic yang tidak disengaja
- Menggunakan interface untuk minimal 1 abstraksi nyata

---

### FASE 3: Concurrency — Jantung Golang (Minggu 5-6)

**Tujuan:** Menguasai model konkurensi Go yang jadi pembeda utama dengan bahasa lain.

**Materi rinci:**
- Goroutine: `go func(){}()`, cara kerja scheduler Go (M:N threading, high level)
- Channel: unbuffered vs buffered, directional channel (`chan<-`, `<-chan`)
- Channel closing, deteksi channel closed (`v, ok := <-ch`)
- `select` statement, `default` case, timeout pattern
- `sync.WaitGroup` untuk menunggu goroutine selesai
- `sync.Mutex` dan `sync.RWMutex` untuk melindungi shared state
- `sync.Once`, `sync.Atomic` (operasi atomic sederhana)
- Race condition: cara terjadi, cara mendeteksi (`go run -race`, `go test -race`)
- Deadlock: penyebab umum dan cara menghindari
- `context.Context`: cancellation, timeout (`context.WithTimeout`), value propagation
- Pola konkurensi umum: worker pool, fan-in/fan-out, pipeline
- Buffered channel sebagai semaphore sederhana

**Deliverable:**
- Project "Worker Pool Pemrosesan Data" — simulasi banyak job diproses paralel dengan limit worker
- Project "Web Scraper Paralel" — fetch banyak URL sekaligus dengan context timeout
- Eksperimen sengaja membuat race condition, lalu perbaiki dengan mutex/channel

**Definition of Done:**
- Semua project lulus `go run -race` tanpa warning
- Bisa menjelaskan perbedaan channel vs mutex dan kapan pakai yang mana

---

### FASE 4: Tooling, Testing, dan Project Structure (Minggu 7)

**Tujuan:** Membiasakan diri dengan praktik profesional dalam menulis kode Go.

**Materi rinci:**
- Go Modules: `go.mod`, `go.sum`, `go mod init/tidy/vendor`
- Struktur project idiomatic: `/cmd`, `/internal`, `/pkg`, `/api`, `/config`
- Testing dengan package `testing`: `TestXxx`, `t.Run` untuk subtest
- Table-driven test (pola standar di Go)
- `go test -cover`, coverage report
- Mocking (manual mock via interface, atau pakai `gomock`/`testify`)
- Benchmarking: `func BenchmarkXxx(b *testing.B)`
- Linting: `golangci-lint`, `go vet`, `gofmt`/`goimports`
- Logging terstruktur (`log/slog` bawaan Go, atau `zap`/`zerolog`)
- Konfigurasi aplikasi: environment variable, `viper`, file `.env`
- Dependency injection sederhana (manual, tanpa framework berat)

**Deliverable:**
- Menambahkan unit test & benchmark ke semua project Fase 1-3
- Setup linter di project (`.golangci.yml`)

**Definition of Done:**
- Coverage test minimal 60% pada logic inti
- Semua kode lolos linter tanpa warning signifikan

---

### FASE 5: REST API Development (Minggu 8-10)

**Tujuan:** Mampu membangun REST API yang production-ready.

**Materi rinci:**
- HTTP dasar di Go: `net/http`, `http.HandleFunc`, `ResponseWriter`, `Request`
- Framework web: pilih salah satu (Gin **direkomendasikan untuk mulai**, atau Fiber, Echo)
- Routing: path parameter, query parameter, route grouping
- Middleware: logging, auth, CORS, rate limiting, recovery dari panic
- Request validation (`go-playground/validator`)
- JSON encoding/decoding (`encoding/json`, struct tag `json:"..."`)
- Autentikasi: JWT (generate & verify token), session-based sebagai pembanding
- Autorisasi: role-based access control sederhana
- Error response yang konsisten (format API error standar)
- API documentation: Swagger/OpenAPI (`swaggo`)
- Pagination, filtering, sorting pada endpoint list
- File upload handling
- Rate limiting & throttling dasar
- Health check endpoint

**Deliverable:**
- Project besar: **"Mini E-Commerce API"**
  - Endpoint: register/login (JWT), CRUD produk, keranjang belanja, checkout
  - Middleware auth & logging
  - Dokumentasi Swagger
  - Unit test untuk handler & service layer

**Definition of Done:**
- API bisa diuji end-to-end via Postman/curl
- Autentikasi & otorisasi berfungsi benar
- Error handling konsisten di semua endpoint

---

### FASE 6: Database — SQL & NoSQL (Minggu 11-12)

**Tujuan:** Mengintegrasikan aplikasi Go dengan database secara aman dan efisien.

**Materi rinci — SQL (PostgreSQL/MySQL):**
- `database/sql` package, driver (`pgx`, `mysql driver`)
- Connection pooling (`SetMaxOpenConns`, `SetMaxIdleConns`)
- Query dasar: `Query`, `QueryRow`, `Exec`
- Prepared statement, mencegah SQL Injection
- Transaction: `Begin`, `Commit`, `Rollback`, isolation level
- ORM: GORM (atau query builder `sqlx`) — trade-off ORM vs raw SQL
- Migration (`golang-migrate` atau `goose`)
- Indexing dasar & query optimization (`EXPLAIN ANALYZE`)
- N+1 query problem dan solusinya

**Materi rinci — NoSQL:**
- Redis: caching pattern (cache-aside), session storage, rate limiting dengan Redis, pub/sub
- MongoDB dasar (opsional, tergantung kebutuhan): koneksi via driver resmi, CRUD, aggregation dasar

**Deliverable:**
- Integrasi database PostgreSQL ke "Mini E-Commerce API" (ganti in-memory storage sebelumnya)
- Tambahkan Redis untuk caching daftar produk & rate limiting login
- Migration script untuk skema database

**Definition of Done:**
- Semua operasi CRUD tersimpan persisten di database
- Transaction pada proses checkout aman dari race condition (stok tidak minus)
- Caching terbukti mengurangi query ke database (bisa dibuktikan dengan log/benchmark)

---

### FASE 7: Microservices & Payment System Pattern (Minggu 13-15)

**Tujuan:** Memahami arsitektur terdistribusi dan pola khusus sistem pembayaran — ini yang membedakan kandidat middle dari senior.

**Materi rinci — Microservices:**
- Kapan microservices dibutuhkan (dan kapan tidak — trade-off vs monolith)
- Komunikasi antar service: REST vs gRPC (Protocol Buffers)
- Service discovery dasar
- API Gateway (konsep)
- Circuit breaker pattern (`sony/gobreaker` atau manual)
- Containerization: Dockerfile untuk Go app (multi-stage build), Docker Compose untuk orkestrasi lokal
- Message queue: Kafka atau RabbitMQ — konsep producer/consumer, event-driven architecture
- Distributed logging & tracing dasar (correlation ID antar service)

**Materi rinci — Payment System Pattern (kritis untuk role ini):**
- **Idempotency**: idempotency key, penyimpanan status request untuk cegah double processing
- **Distributed transaction**: 2-phase commit (konsep), **Saga pattern** (choreography vs orchestration)
- **Eventual consistency** vs strong consistency — trade-off di sistem pembayaran
- **Reconciliation**: pencocokan data transaksi antara sistem internal dan payment gateway
- **Webhook handling** dari payment gateway (Midtrans/Xendit) — verifikasi signature, retry mechanism
- **Race condition pada saldo/stok**: optimistic locking vs pessimistic locking
- **Outbox pattern** untuk menjamin konsistensi antara database dan message queue
- Audit trail & logging transaksi finansial

**Deliverable:**
- Pecah "Mini E-Commerce API" menjadi 2 service:
  - **Order Service**: kelola pesanan, panggil Payment Service via gRPC
  - **Payment Service**: proses pembayaran dengan idempotency key, simulasi callback webhook
- Implementasi saga pattern sederhana: jika pembayaran gagal, order otomatis di-cancel (compensating transaction)
- Docker Compose untuk menjalankan kedua service + database + Redis sekaligus

**Definition of Done:**
- Uji coba double-submit request pembayaran — hasil tetap konsisten (idempotency terbukti)
- Simulasi kegagalan payment service — order ter-cancel otomatis (saga compensating action jalan)
- Seluruh sistem bisa dijalankan dengan `docker-compose up`

---

### FASE 8: Observability, Deployment, & Persiapan Interview (Minggu 16)

**Tujuan:** Melengkapi kemampuan operasional dan siap menghadapi proses rekrutmen.

**Materi rinci:**
- Monitoring dasar: metrics dengan Prometheus (expose `/metrics`), visualisasi Grafana (konsep)
- Health check & readiness/liveness probe (relevan jika nanti masuk Kubernetes)
- CI/CD sederhana: GitHub Actions untuk `go test` & `go build` otomatis
- Review menyeluruh seluruh materi Fase 1-7
- Latihan soal algoritma umum dalam Go (slice manipulation, string processing, concurrency puzzle)
- Latihan system design level dasar-menengah (desain sistem checkout, desain rate limiter, desain URL shortener)
- Mock interview: technical question (konsep Go) + live coding + system design

**Deliverable:**
- Portofolio GitHub rapi berisi seluruh project (fundamentals, inventaris, worker pool, mini e-commerce, microservices payment)
- README di setiap project menjelaskan cara menjalankan & keputusan desain
- Catatan ringkasan konsep (cheat sheet pribadi) untuk review cepat sebelum interview

**Definition of Done:**
- Lulus simulasi mock interview (bisa dengan bantuan AI atau rekan)
- Bisa menjelaskan setiap keputusan desain di project payment system

---

## 5. Ringkasan Timeline

```
Fase 1  (Minggu 1-2)   : Fundamental Syntax
Fase 2  (Minggu 3-4)   : Struct, Interface, Error Handling, Pointer
Fase 3  (Minggu 5-6)   : Concurrency
Fase 4  (Minggu 7)     : Tooling, Testing, Project Structure
Fase 5  (Minggu 8-10)  : REST API Development
Fase 6  (Minggu 11-12) : Database SQL & NoSQL
Fase 7  (Minggu 13-15) : Microservices & Payment Pattern
Fase 8  (Minggu 16)    : Observability, Deployment, Persiapan Interview
```

## 6. Functional Requirements (Kebutuhan Belajar)

| ID | Kebutuhan | Prioritas |
|----|-----------|-----------|
| FR-1 | Sumber belajar resmi (go.dev/tour, go.dev/doc) sebagai referensi utama | Must |
| FR-2 | Environment coding lengkap (Go, Docker, PostgreSQL, Redis lokal) | Must |
| FR-3 | Setiap fase menghasilkan project nyata, bukan sekadar latihan soal | Must |
| FR-4 | Kode disimpan di GitHub sebagai portofolio | Must |
| FR-5 | Testing & linting diterapkan sejak Fase 4 dan seterusnya | Should |
| FR-6 | Simulasi interview di akhir program | Should |
| FR-7 | Review mingguan progress dibanding rencana | Could |

## 7. Non-Functional Requirements

- **Konsistensi belajar:** minimal 2 jam/hari, 5-6 hari/minggu
- **Rasio teori vs praktik:** ideal 25:75
- **Dokumentasi:** setiap project punya README yang jelas
- **Reproducibility:** semua project bisa dijalankan orang lain hanya dengan mengikuti README

## 8. Risiko & Mitigasi

| Risiko | Dampak | Mitigasi |
|--------|--------|----------|
| Overwhelmed karena materi banyak | Motivasi turun | Ikuti urutan fase, jangan loncat, pecah jadi target harian kecil |
| Konsep concurrency terasa sulit | Fase 3 molor | Alokasikan waktu ekstra, banyak latihan program kecil sebelum project besar |
| Payment pattern terasa abstrak tanpa konteks industri | Sulit dijelaskan saat interview | Baca studi kasus nyata (engineering blog Tokopedia, Shopee, Gojek, Stripe) |
| Terlalu fokus tutorial, kurang menulis kode sendiri | Skill tidak melekat | Wajib deliverable project di setiap fase, bukan hanya mengikuti tutorial |
| Kesulitan database transaction & race condition | Bug sulit dideteksi | Sengaja buat test case concurrent (multiple goroutine checkout produk stok terbatas) |

## 9. Kriteria Keberhasilan Akhir Program

- [ ] Minimal 5 project di GitHub: Fundamentals, Manajemen Inventaris, Worker Pool/Scraper, Mini E-Commerce API, Microservices Payment System
- [ ] Semua project punya unit test dengan coverage layak (≥60% pada logic inti)
- [ ] Mampu menjelaskan trade-off teknis (mutex vs channel, SQL vs NoSQL, monolith vs microservices, ORM vs raw SQL)
- [ ] Mampu mengimplementasikan idempotency dan saga pattern dari nol
- [ ] Lulus simulasi mock interview teknis level middle-senior
- [ ] Siap melamar dan mengikuti technical test untuk posisi Golang Developer e-commerce/payment