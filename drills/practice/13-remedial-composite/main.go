package main

import (
	"encoding/json"
	"fmt"
)

// ============================================================================
// SPESIFIKASI DATA MODEL (FINTECH DIGITAL WALLET)
// ============================================================================

// 1. Struct Currency
// - Code: string, tag json: "code" (contoh: "USD", "SGD", "JPY")
// - Name: string, tag json: "name" (contoh: "US Dollar")
// - RateToIDR: int, tag json: "rate_to_idr" (kurs ke Rupiah)
type Currency struct {
	Code      string `json:"code"`
	Name      string `json:"name"`
	RateToIDR int    `json:"rate_to_idr"`
}

// 2. Struct BalanceItem (Wajib gunakan Struct Embedding!)
// - Embed struct Currency (Anonymous field)
// - Amount: int, tag json: "amount"
type BalanceItem struct {
	// TODO 1: Embed struct Currency di sini
	// TODO 2: Tambahkan field Amount
	Currency
	Amount int `json:"amount"`
}

// 3. Struct Wallet
// - OwnerID: string, tag json: "owner_id"
// - Balances: []BalanceItem, tag json: "balances"
// - InternalPinHash: string, tag json: "-" (RAHASIA, tidak boleh keluar di JSON!)
type Wallet struct {
	OwnerID         string        `json:"owner_id"`
	Balances        []BalanceItem `json:"balances"`
	InternalPinHash string        `json:"-"`
}

// ============================================================================
// LOGIKA BISNIS (TULIS SENDIRI TANPA COPAS/BANTUAN AI)
// ============================================================================

// TODO 3: Implementasikan Deposit
// Aturan:
//  1. Cek apakah currency code terdaftar di map 'supported' (gunakan comma-ok idiom!).
//     Jika tidak terdaftar, return false.
//  2. Periksa apakah mata uang tersebut sudah ada di wallet.Balances:
//     - Jika sudah ada, tambahkan Amount-nya: wallet.Balances[i].Amount += amount
//     - Langsung return true.
//  3. Jika setelah diperiksa seluruhnya ternyata belum ada (atau Balances masih kosong):
//     - Buat BalanceItem baru lalu append ke wallet.Balances.
//     - Return true.
func Deposit(supported map[string]Currency, wallet *Wallet, code string, amount int) bool {
	exCurr, ok := supported[code]
	if !ok {
		return false
	}

	for i, item := range wallet.Balances {
		if item.Currency.Code == code {
			wallet.Balances[i].Amount += amount
			return true
		}
	}

	wallet.Balances = append(wallet.Balances, BalanceItem{
		Currency: exCurr,
		Amount:   amount,
	})
	return true
}

// TODO 4: Implementasikan CloseCurrency
// Aturan:
//  1. Cari BalanceItem dengan Code yang sesuai di wallet.Balances.
//  2. Jika ketemu, hapus elemen tersebut dari slice wallet.Balances.
//     (Gunakan teknik potong & sambung slice yang sudah kamu kuasai).
//     Return true.
//  3. Jika tidak ditemukan, return false.
func CloseCurrency(wallet *Wallet, code string) bool {
	for i, item := range wallet.Balances {
		if item.Currency.Code == code {
			wallet.Balances = append(wallet.Balances[:i], wallet.Balances[i+1:]...)
			return true
		}
	}
	return false
}

// TODO 5: Implementasikan TotalPortfolioIDR
// Hitung total seluruh saldo dalam Rupiah:
// Jumlahkan (Amount * RateToIDR) dari semua item di wallet.Balances.
func TotalPortfolioIDR(wallet Wallet) int {
	total := 0
	for _, item := range wallet.Balances {
		total += item.Amount * item.Currency.RateToIDR
	}
	return total
}

func main() {
	fmt.Println("=== REMEDIAL INDEPENDENT DRILL: MULTI-CURRENCY WALLET ===")

	// 1. Master Data Mata Uang yang Didukung (Map)
	supportedCurrencies := map[string]Currency{
		"USD": {Code: "USD", Name: "US Dollar", RateToIDR: 16000},
		"SGD": {Code: "SGD", Name: "Singapore Dollar", RateToIDR: 12000},
		"JPY": {Code: "JPY", Name: "Japanese Yen", RateToIDR: 100},
	}

	// 2. Inisialisasi Wallet
	wallet := Wallet{
		OwnerID:         "CUST-007",
		Balances:        make([]BalanceItem, 0, 4),
		InternalPinHash: "$2a$12$e8YQe9r8H3Fj2kL8...",
	}

	// 3. Uji Deposit
	fmt.Println("1. Deposit USD 100:", Deposit(supportedCurrencies, &wallet, "USD", 100))                // Harusnya true
	fmt.Println("2. Deposit SGD 200:", Deposit(supportedCurrencies, &wallet, "SGD", 200))                // Harusnya true
	fmt.Println("3. Deposit USD 50 (akumulasi):", Deposit(supportedCurrencies, &wallet, "USD", 50))      // Harusnya true (USD jadi 150)
	fmt.Println("4. Deposit BITCOIN (tidak didukung):", Deposit(supportedCurrencies, &wallet, "BTC", 1)) // Harusnya false

	// 4. Uji TotalPortfolioIDR
	// Perhitungan:
	// - USD: 150 * 16.000 = Rp 2.400.000
	// - SGD: 200 * 12.000 = Rp 2.400.000
	// Total: Rp 4.800.000
	totalIDR := TotalPortfolioIDR(wallet)
	fmt.Printf("5. Total Portfolio: Rp %d (Ekspektasi: Rp 4800000)\n", totalIDR)

	// 5. Uji CloseCurrency
	fmt.Println("6. Hapus/Tutup Akun SGD:", CloseCurrency(&wallet, "SGD")) // Harusnya true
	totalSetelahTutup := TotalPortfolioIDR(wallet)
	// Sisa USD saja: 150 * 16.000 = Rp 2.400.000
	fmt.Printf("7. Total Setelah Tutup SGD: Rp %d (Ekspektasi: Rp 2400000)\n", totalSetelahTutup)

	// 6. Uji Serialisasi JSON (Pastikan InternalPinHash TIDAK MUNCUL!)
	jsonBytes, _ := json.Marshal(wallet)
	fmt.Println("8. Output JSON Wallet:")
	fmt.Println(string(jsonBytes))
}
