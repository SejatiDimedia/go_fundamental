package main

import (
	"fmt"
)

// =========================================================================
// MINI-CHALLENGE TOPIK 1: FINANCIAL DATA MASKER & FORMATTER
// Menggabungkan: CLI args, Enum iota, Zero-value defense, int64 money, Runes
// =========================================================================

// TODO 1: Buat custom type 'AccountType int' dan konstanta enum via iota
// - AccUnknown = 0 (Zero value defense)
// - AccSavings = 1
// - AccChecking = 2
// const ( ... )
type AccountType int

const (
	AccUnknown AccountType = iota
	AccSavings
	AccChecking
)

// TODO 2: Buat fungsi MaskAccountNumber(accNumber string) string
//   - Ubah string ke []rune
//   - Jika panjang karakter <= 6, kembalikan text aslinya (terlalu pendek untuk dimask)
//   - Jika panjang > 6, sensor digit tengah: pertahankan 3 karakter pertama dan 3 karakter terakhir,
//     ganti karakter di antaranya dengan '*'
//   - Contoh: "1234567890" -> "123****890"
//
// func MaskAccountNumber ...
func MaskAccountNumber(accNumber string) string {
	runes := []rune(accNumber)
	if len(runes) <= 6 {
		return accNumber
	}

	for i := 3; i < len(runes)-3; i++ {
		runes[i] = '*'
	}
	return string(runes)
}

// TODO 3: Buat fungsi FormatCurrency(amountInSmallestUnit int64, currencySymbol string) string
// - Menggunakan fmt.Sprintf untuk format string
// - Contoh input: 500000, "Rp" -> return "Rp 500000"
// func FormatCurrency ...
func FormatCurrency(amountInSmallestUnit int64, currencySymbol string) string {
	return fmt.Sprintf("%s %d", currencySymbol, amountInSmallestUnit)
}

func main() {
	// TODO 4:
	// a. Uji AccountType: buat variabel tanpa inisialisasi (zero value), cetak nilainya (harus 0)
	var accType AccountType
	fmt.Println(accType)
	// b. Uji MaskAccountNumber:
	//    - MaskAccountNumber("08123456789") -> ekspektasi: "081*****789"
	fmt.Println(MaskAccountNumber("08123456789"))
	// c. Uji FormatCurrency:
	//    - FormatCurrency(250000, "Rp") -> ekspektasi: "Rp 250000"
	fmt.Println(FormatCurrency(250000, "Rp"))

}
