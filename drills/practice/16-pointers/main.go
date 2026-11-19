package main

import "fmt"

type Account struct {
	ID      string
	Balance int
}

// TODO 1: Implementasikan fungsi Swap(a, b *int)
// Fungsi ini menukar nilai dari dua variabel integer di pemanggil
// Contoh: jika a=10, b=20, setelah Swap dipanggil maka a=20, b=10
// Petunjuk: gunakan dereference (*) untuk membaca & menukar nilainya
func Swap(a, b *int) {
	// Tulis logika tukar di sini:
	tmp := *a
	*a = *b
	*b = tmp
}

// TODO 2: Implementasikan fungsi DeductBalance(acc *Account, amount int) bool
// - Periksa apakah acc.Balance >= amount
// - Jika cukup, kurangi acc.Balance dengan amount, lalu return true
// - Jika tidak cukup, jangan ubah saldo, dan return false
func DeductBalance(acc *Account, amount int) bool {
	if acc.Balance >= amount {
		acc.Balance -= amount
		return true
	}
	return false
}

// TODO 3: Implementasikan fungsi NewAccount(id string, initialBalance int) *Account
// - Ini adalah pola Constructor Pattern di Go
// - Buat struct Account baru dengan nilai id dan initialBalance
// - Return POINTER (&) ke struct Account tersebut
func NewAccount(id string, initialBalance int) *Account {
	return &Account{
		ID:      id,
		Balance: initialBalance,
	}
}

func main() {
	fmt.Println("=== 1. UJI POINTER SWAP ===")
	x, y := 10, 20
	fmt.Printf("Sebelum Swap: x=%d, y=%d\n", x, y)
	Swap(&x, &y)
	fmt.Printf("Setelah Swap : x=%d, y=%d (Ekspektasi: x=20, y=10)\n", x, y)

	fmt.Println("\n=== 2. UJI CONSTRUCTOR PATTERN ===")
	myAcc := NewAccount("ACC-001", 100_000)
	if myAcc == nil {
		fmt.Println("myAcc masih nil! Lengkapi TODO 3.")
		return
	}
	fmt.Printf("Akun Baru: ID=%s, Saldo=%d (Alamat: %p)\n", myAcc.ID, myAcc.Balance, myAcc)

	fmt.Println("\n=== 3. UJI MUTASI STRUCT VIA POINTER ===")
	// Coba potong saldo 40.000 (harus berhasil)
	success := DeductBalance(myAcc, 40_000)
	fmt.Printf("Potong 40.000: %t, Sisa Saldo: %d (Ekspektasi: 60000)\n", success, myAcc.Balance)

	// Coba potong saldo 90.000 (harus gagal karena saldo sisa 60.000)
	failed := DeductBalance(myAcc, 90_000)
	fmt.Printf("Potong 90.000: %t, Sisa Saldo: %d (Ekspektasi: 60000)\n", failed, myAcc.Balance)
}
