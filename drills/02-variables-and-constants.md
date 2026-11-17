# Catatan Deep Dive: Variables, Constants, & Scope

Dokumen ini merangkum seluruh hasil diskusi dan deep dive pada sub-topik **1. Language Basics > Variables & Constants**, mencakup sintaks, zero values, shadowing, dan konstanta enum.

---

## 1. `var` vs `:=` (Short Variable Declaration)

| Aturan | `var` | `:=` (Short Declaration) |
| :--- | :--- | :--- |
| **Lingkup (Scope)** | Bisa di level package (luar fungsi) maupun di dalam fungsi. | **HANYA bisa di dalam fungsi** (block scope). Dilarang keras di luar fungsi. |
| **Kapan Digunakan** | 1. Ingin memanfaatkan **Zero Value** (misal `var count int`).<br>2. Di package-level scope.<br>3. Ingin tipe eksplisit (misal `var id int64 = 1`). | Ketika variabel lokal memiliki nilai awal konkret atau menampung return fungsi: `user, err := getUser()`. |

---

## 2. Konsep Zero Values di Go

Di Go, tidak ada konsep *uninitialized garbage memory*. Semua variabel yang dideklarasikan tanpa nilai awal otomatis diberi nilai dasar yang aman:
* `int`, `float`: `0`, `0.0`
* `bool`: `false`
* `string`: `""` (empty string)
* `pointer`, `slice`, `map`, `channel`, `interface`: **`nil`**
* `struct`: Setiap field di dalam struct otomatis bernilai zero value masing-masing.

### Perilaku `nil slice` vs `nil map` (Crucial Gotcha):
1. **Slice:**  
   `var s []int` (nil).  
   `s = append(s, 10)` 👉 **AMAN dan SANGAT IDIOMATIC.** `append` otomatis mengalokasikan array baru.
2. **Map:**  
   `var m map[string]int` (nil).  
   * Membaca: `val := m["key"]` 👉 **AMAN** (mengembalikan 0).  
   * Menulis: `m["key"] = 10` 👉 **PANIC!** (`panic: assignment to entry in nil map`).  
   * Solusi: Wajib diinisialisasi dengan `make(map[string]int)` sebelum ditulis.

---

## 3. Variable Shadowing (Silent Killer di Production)

Shadowing terjadi saat variabel dideklarasikan ulang dengan nama yang sama di scope yang lebih sempit (misal di dalam blok `if` atau `for`), sehingga variabel luar tertutup.

```go
var user *User
var err error

if condition {
    // SALAH: menggunakan := membuat user baru di dalam if
    user, err := getUser() 
    if err != nil {
        return err
    }
}
// Di luar if, 'user' MASIH nil!
println(user.Name) // PANIC: nil pointer dereference
```

**Solusi:**  
Deklarasikan variabel `err` di luar, lalu gunakan operator penugasan biasa (**`=`**), bukan `:=`:
```go
user, err = getUser()
```

---

## 4. `const` dan `iota`

1. **`iota` untuk Enum:**  
   `iota` adalah pencacah otomatis yang dimulai dari angka `0` dan bertambah `+1` pada setiap baris dalam blok `const`.
2. **Zero Value Defense:**  
   Jangan gunakan nilai `0` untuk status valid (seperti `Approved` atau `Admin`). Buat nilai `0` menjadi `Unknown` atau buang menggunakan blank identifier (`_ = iota`) agar saat struct tidak sengaja diinisialisasi dengan zero value, sistem tidak memberikan akses tidak sah.
3. **Kenapa `const` Hanya untuk Basic Types?**  
   * `const` harus bernilai mutlak saat **Compile-Time** (di-hardcode ke binary).
   * `slice`, `map`, dan `struct` membutuhkan alokasi memori RAM yang baru terjadi saat **Runtime**.
   * Elemen internal `slice` dan `map` bersifat *mutable*, sedangkan `const` menjamin *immutability*.
