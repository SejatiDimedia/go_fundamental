package main

import (
	"fmt"
	"sync"
)

type BankAccount struct {
	mu      sync.Mutex
	Balance int
}

// SALAH: Value Receiver memfotokopi Mutex!
func (b BankAccount) DepositValue(amount int) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.Balance += amount
}

// BENAR: Pointer Receiver menjaga Mutex yang sama
func (b *BankAccount) DepositPointer(amount int) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.Balance += amount
}

func main() {
	acc := BankAccount{Balance: 100}
	acc.DepositValue(50)
	fmt.Println("Balance:", acc.Balance)
}
