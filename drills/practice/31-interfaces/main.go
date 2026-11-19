package main

import "fmt"

type Payer interface {
	Pay(amount int)
}

type BankTransfer struct {
	Balance int
}

// Method dengan Pointer Receiver
func (b *BankTransfer) Pay(amount int) {
	b.Balance -= amount
}

func ProcessPayment(p Payer, amount int) {
	p.Pay(amount)
}

func main() {
	// Skenario 1: Mengoper VALUE
	// bVal := BankTransfer{Balance: 100_000}
	// ProcessPayment(bVal, 50_000) // ❌ COMPILE ERROR!

	// Skenario 2: Mengoper POINTER
	bPtr := &BankTransfer{Balance: 100_000}
	ProcessPayment(bPtr, 50_000) // ✅ SUKSES!

	fmt.Println("Sisa saldo:", bPtr.Balance)
}
