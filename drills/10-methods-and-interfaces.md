# Catatan Deep Dive: Methods and Interfaces di Go

Dokumen ini merangkum seluruh konsep, cara kerja memori, jebakan produksi, dan best practice pada **Topik 6: Methods and Interfaces**.

---

## 1. Methods vs Functions

Method adalah fungsi yang memiliki parameter khusus bernama **Receiver** yang ditempelkan ke sebuah tipe data:

```go
// Sintaks Method:
func (r ReceiverType) MethodName(params) returns { ... }
```

Di balik layar, compiler Go mengubah method menjadi fungsi biasa dengan receiver sebagai parameter pertamanya:
`u.ChangeName("Bob")` $\rightarrow$ `(*User).ChangeName(u, "Bob")`.

---

## 2. Value Receiver vs Pointer Receiver

| Karakteristik | Value Receiver `(u User)` | Pointer Receiver `(u *User)` |
| :--- | :--- | :--- |
| **Operasi Memori** | Meng-copy seluruh struct (Pass by Value) | Mengoper alamat memori asli struct |
| **Mutasi Data** | **TIDAK BISA** mengubah data asli pemanggil | **BISA** memodifikasi data asli pemanggil |
| **Nil Receiver** | Panic jika dipanggil pada nil | **BISA** menerima `nil` secara aman jika di-handle |
| **Penggunaan Ideal** | Tipe kecil ($\le 64$ byte), immutable/read-only | Butuh mutasi, struct besar, atau ada `sync.Mutex` |

### ⚠️ The Map Method Trap
Elemen di dalam map **tidak addressable**. Maka memanggil pointer method pada elemen map (`users["admin"].ChangeName("Bob")`) akan **compile error**!
* **Solusi:** Simpan pointer di dalam map: `map[string]*User`.

### 🔒 Concurrency Safety (`sync.Mutex`)
Jika sebuah struct memiliki field `sync.Mutex`, **SELURUH METHOD-NYA WAJIB MENGGUNAKAN POINTER RECEIVER**.
* Menggunakan value receiver akan memfotokopi gembok mutex (*passes lock by value*), menyebabkan data race / lock yang rusak antar-goroutine.

---

## 3. Interfaces & Implicit Satisfaction (Duck Typing)

Go **tidak memiliki kata kunci `implements`**.
Sebuah tipe data secara otomatis memenuhi (*satisfy*) interface jika tipe tersebut memiliki seluruh method dengan nama, parameter, dan return type yang sama persis seperti yang dideklarasikan oleh interface.

> *"If it walks like a duck, swims like a duck, and quacks like a duck, it's a duck."*

---

## 4. The Method Sets Rule (`T` vs `*T` pada Interface)

Ketika memasukkan nilai konkret ke dalam variabel Interface:

| Tipe Nilai yang Dimasukkan | Method dengan `(t T)` | Method dengan `(t *T)` |
| :--- | :---: | :---: |
| **Value Biasa (`T`)** | ✅ Memenuhi Interface | ❌ **GAGAL (Compile Error)** |
| **Pointer (`*T`)** | ✅ Memenuhi Interface | ✅ **Memenuhi Interface** |

* Data di dalam interface bersifat *unaddressable*. Oleh karena itu, jika sebuah method menggunakan **Pointer Receiver**, kamu **wajib memasukkan pointer `&MyStruct{}`** ke dalam interface.

---

## 5. Empty Interface (`any` / `interface{}`) & Type Assertion

`any` (alias dari `interface{}`) memiliki **0 buah method**, sehingga **seluruh tipe data di Go otomatis memenuhi `any`**.

### A. Type Assertion & The Comma-Ok Idiom
Type assertion **BUKAN** konversi tipe, melainkan pemeriksaan tipe konkret di dalam interface:
* **Tanpa Comma-Ok:** `val := data.(int)` $\rightarrow$ **PANIC** jika tipe salah!
* **Dengan Comma-Ok (Best Practice):** `val, ok := data.(int)` $\rightarrow$ Aman, `ok = false` jika salah tipe.

### B. Type Switch
Pola elegan untuk menangani berbagai tipe data dalam satu blok kode:
```go
switch v := data.(type) {
case int:
    fmt.Println("Integer:", v + 10)
case string:
    fmt.Println("String:", len(v))
default:
    fmt.Printf("Unknown: %T\n", v)
}
```

---

## 6. Interface Embedding & The Golden Rule

### A. Interface Embedding (Composition)
Menggabungkan beberapa interface kecil menjadi satu interface baru:
```go
type Reader interface { Read() string }
type Writer interface { Write(string) }

type ReadWriter interface {
    Reader
    Writer
}
```

### B. The Golden Rule of Go Architecture
> **"Accept Interfaces, Return Structs"**
> * Terima interface pada parameter fungsi (membuat fungsi fleksibel & mudah di-mock pada unit test).
> * Kembalikan concrete struct pada return value fungsi (memberikan fleksibilitas penuh kepada pemanggil).
