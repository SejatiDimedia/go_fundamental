package main

import "fmt"

type Wallet struct {
	Balance int
}

// Perhatikan posisi bintang: (w *Wallet), bukan (*w Wallet)
func (w *Wallet) Withdraw(amount int) bool {
	if w.Balance >= amount {
		w.Balance -= amount // Kurangi saldo!
		return true
	}
	return false
}

func main() {
	myWallet := Wallet{Balance: 100_000}
	fmt.Println("Saldo awal:", myWallet.Balance)

	// Tarik 40.000 (sukses)
	success := myWallet.Withdraw(40_000)
	fmt.Printf("Tarik 40.000: Sukses? %t, Sisa Saldo: %d\n", success, myWallet.Balance)

	// Tarik 80.000 (gagal, saldo tidak cukup)
	failed := myWallet.Withdraw(80_000)
	fmt.Printf("Tarik 80.000: Sukses? %t, Sisa Saldo: %d\n", failed, myWallet.Balance)
}
