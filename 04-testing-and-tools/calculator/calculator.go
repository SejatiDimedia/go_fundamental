package calculator

import (
	"errors"
	"fmt"
)

var (
	ErrDivisionByZero = errors.New("tidak bisa membagi dengan angka nol")
	ErrNegativePrice  = errors.New("harga tidak boleh negatif")
)

// CalculateTotalDiscount menghitung total harga setelah diskon
func CalculateTotalDiscount(price float64, discountPercent int) (float64, error) {
	if price < 0 {
		return 0, ErrNegativePrice
	}
	if discountPercent < 0 || discountPercent > 100 {
		return 0, fmt.Errorf("persentase diskon %d tidak valid (harus 0-100)", discountPercent)
	}

	discountAmount := (price * float64(discountPercent)) / 100.0
	return price - discountAmount, nil
}

// Divide membagi dua angka dan mencegah pembagian dengan nol
func Divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, ErrDivisionByZero
	}
	return a / b, nil
}
