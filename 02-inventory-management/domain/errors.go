package domain

import "errors"

var (
	ErrProductNotFound     = errors.New("Product tidak di temukandi inventaris")
	ErrOutOfStock          = errors.New("Stock product tidak mencukupi")
	ErrInsufficientPayment = errors.New("Nominal pembayaran kurang dari total belanja")
	ErrInvalidQuantity     = errors.New("Jumlah barang harus lebih dari 0")
)
