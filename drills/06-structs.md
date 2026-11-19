# Catatan Deep Dive: Structs di Go

Dokumen ini merangkum seluruh hasil diskusi dan deep dive pada sub-topik **2. Composite Types > Structs**, mencakup aturan visibility (exported/unexported), tag JSON (`omitempty`, `-`), struct embedding (composition over inheritance), serta trade-off pass-by-value vs pass-by-pointer.

---

## 1. Exported vs Unexported (Visibility Rule)

Di Go, hak akses ditentukan semata-mata oleh **huruf pertama nama identifier**:
* **Huruf Besar (Exported / Public):** Bisa diakses oleh package lain.
* **Huruf Kecil (Unexported / Private):** Hanya bisa diakses di dalam package yang sama.

### 🚨 Jebakan Fatal pada `encoding/json`:
Package `encoding/json` adalah package luar yang membaca struct menggunakan *reflection*.
* Field berawalan **huruf kecil AKAN DIABAIKAN TOTAL** saat `Marshal` maupun `Unmarshal`.
* **Solusi:** Selalu gunakan huruf kapital untuk field yang akan di-serialize ke JSON (misal `Email string` bukan `email string`).

---

## 2. Anatomi Struct Tag & Aturan Khusus

Format penulisan tag: `` `key:"value,option"` ``

### Fitur-Fitur Tag JSON Populer:
1. **Mapping Nama Field:** `Name string `json:"full_name"``
2. **Mengabaikan Field Rahasia:** `PasswordHash string `json:"-"``
   * Field ini tidak akan pernah keluar di JSON output, dan tidak akan bisa diisi lewat JSON input.
3. **Opsi `omitempty`:**
   * Membuang field dari output JSON jika nilainya adalah **zero value** (string `""`, int `0`, pointer `nil`, bool `false`).
   * ⚠️ **Awas Jebakan Angka 0:** Pada data finansial (seperti `Balance int`), jika nilai saldo adalah `0`, tag `omitempty` akan menghapus field `balance` dari response! Frontend tidak akan menerima key `balance`. Jangan gunakan `omitempty` jika angka `0` adalah nilai bisnis yang sah.

---

## 3. Composition over Inheritance (Struct Embedding)

Go menolak konsep pewarisan kelas hierarkis (`class A extends B`) untuk menghindari:
1. *Brittle Base Class Problem* (perubahan di class induk merusak class anak).
2. *Gorilla and Banana Problem* (terpaksa mewarisi seluruh kompleksitas class induk).

### Cara Kerja Struct Embedding:
Go menggunakan *Anonymous Field* di dalam struct:

```go
type AuditLog struct {
    CreatedAt string `json:"created_at"`
    UpdatedAt string `json:"updated_at"`
}

type Transaction struct {
    ID     string `json:"id"`
    Amount int    `json:"amount"`
    AuditLog      // Anonymous Field / Embedded
}
```

### Field Promotion:
Semua field milik `AuditLog` otomatis dipromosikan ke `Transaction`:
* Bisa diakses langsung: `trx.CreatedAt = "2026-11-18"`
* Saat di-`json.Marshal`, seluruh field anak otomatis rata (*flattened*) di root JSON object:
  ```json
  {"id":"TRX-01","amount":50000,"created_at":"2026-11-18","updated_at":"2026-11-18"}
  ```

### Field Shadowing (Tabrakan Nama Field):
Jika struct luar dan struct dalam memiliki nama field yang sama, struct luar akan menaungi (*shadow*) struct dalam:
* `u.City` mengakses field milik struct terluar.
* Field dalam tetap aman dan bisa diakses via path eksplisit: `u.Address.City`.

---

## 4. Pass by Value vs Pass by Pointer pada Struct

| Aspek | Pass by Value (`func F(u User)`) | Pass by Pointer (`func F(u *User)`) |
| :--- | :--- | :--- |
| **Mutasi** | Read-only. Struct asli pemanggil dijamin aman tidak berubah. | Read-Write. Perubahan di dalam fungsi mengubah struct asli pemanggil. |
| **Ukuran Copy** | Meng-copy seluruh byte struct. Boros jika struct besar. | Selalu 8 byte (alamat memori), sangat efisien untuk struct besar. |
| **Alokasi Memori** | Umumnya dialokasikan di **Stack** (cepat & ramah GC). | Cenderung mengalami *escape to heap* (menambah beban GC). |

### Aturan Praktis:
* Ubah data / model besar / DTO $\rightarrow$ Gunakan **Pointer (`*User`)**.
* Read-only & struct berukuran kecil ($\le$ 64 byte) $\rightarrow$ Gunakan **Value (`User`)**.
