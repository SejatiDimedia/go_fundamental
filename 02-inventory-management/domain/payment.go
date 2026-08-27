package domain

import (
	"fmt"
	"math/rand"
	"time"
)

type PaymentMethod interface {
	Name() string
	Pay(amount int) (string, error)
}

type CashPayment struct {
	CashGiven int
}

func (c *CashPayment) Name() string {
	return "Tunai (Cash)"
}

func (c *CashPayment) Pay(amount int) (string, error) {
	if c.CashGiven < amount {
		return "", fmt.Errorf("uang tunai Rp%d kurang dari total tagihan Rp%d: %w", c.CashGiven, amount, ErrInsufficientPayment)
	}
	change := c.CashGiven - amount
	receipt := fmt.Sprintf("CASH-OK | Kembalian: Rp%d", change)
	return receipt, nil
}

type VAPayment struct {
	BankName string
	VANumber string
}

func (v *VAPayment) Name() string {
	return "Virtual Account " + v.BankName
}

func (v *VAPayment) Pay(amount int) (string, error) {
	refCode := fmt.Sprintf("VA-%s-%d", v.BankName, time.Now().Unix()%1000000)
	return refCode, nil
}

type EWalletPayment struct {
	Provider string
	WalletID string
	Balance  int
}

func (e *EWalletPayment) Name() string {
	return "E-Wallet (" + e.Provider + ")"
}

func (e *EWalletPayment) Pay(amount int) (string, error) {
	if e.Balance < amount {
		return "", fmt.Errorf("saldo %s Rp%d tidak cukup untuk bayar Rp%d: %w", e.Provider, e.Balance, amount, ErrInsufficientPayment)
	}
	e.Balance -= amount
	refCode := fmt.Sprintf("EWL-%s-%04d", e.Provider, rand.Intn(9999))
	return refCode, nil
}
