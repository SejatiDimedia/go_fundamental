package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

// ==========================================================
// 1. REKENING TIDAK AMAN (Sengaja Dibuat Buggy Tanpa Gembok)
// ==========================================================
type UnsafeAccount struct {
	Balance int
}

func (acc *UnsafeAccount) Deposit(amount int) {
	// Terjadi data race di sini jika banyak goroutine memanggil Deposit() bersamaan!
	acc.Balance += amount
}

// ==========================================================
// 2. REKENING AMAN DENGAN MUTEX (sync.Mutex)
// ==========================================================
type MutexAccount struct {
	mu      sync.Mutex // Gembok Mutex
	Balance int
}

func (acc *MutexAccount) Deposit(amount int) {
	acc.mu.Lock()         // 🔒 Kunci gembok: Hanya 1 goroutine yang boleh masuk
	defer acc.mu.Unlock() // 🔓 Buka gembok saat selesai

	acc.Balance += amount
}

// ==========================================================
// 3. REKENING AMAN DENGAN OPERASI ATOMIC (sync/atomic)
// ==========================================================
type AtomicAccount struct {
	Balance int64
}

func (acc *AtomicAccount) Deposit(amount int64) {
	// Operasi penjumlahan di level instruksi CPU (Sangat cepat untuk counter/angka tunggal)
	atomic.AddInt64(&acc.Balance, amount)
}

func main() {
	const totalTransactions = 1000                            // 1.000 Transaksi bersamaan
	const depositAmount = 1000                                // Masing-masing setor Rp1.000
	const expectedBalance = totalTransactions * depositAmount // Ekspektasi Saldo = Rp1.000.000

	fmt.Println("==================================================================")
	fmt.Println("       💥 SIMULASI RACE CONDITION & PEMBUKTIAN MUTEX (Go)         ")
	fmt.Printf("       Total Transaksi : %d transaksi serentak @ Rp%d\n", totalTransactions, depositAmount)
	fmt.Printf("       Ekspektasi Saldo: Rp%d\n", expectedBalance)
	fmt.Println("==================================================================\n")

	// ----------------------------------------------------------
	// Uji Coba 1: Rekening Tidak Aman (Unsafe)
	// ----------------------------------------------------------
	unsafeAcc := &UnsafeAccount{Balance: 0}
	var wg1 sync.WaitGroup

	for i := 0; i < totalTransactions; i++ {
		wg1.Add(1)
		go func() {
			defer wg1.Done()
			unsafeAcc.Deposit(depositAmount)
		}()
	}
	wg1.Wait()

	fmt.Println("--- HASIL EKSPERIMEN 1: TANPA MUTEX (UNSAFE) ---")
	fmt.Printf("Saldo Akhir   : Rp%d\n", unsafeAcc.Balance)
	if unsafeAcc.Balance != expectedBalance {
		fmt.Printf("❌ SALDO KORUP / HILANG: Rp%d hilang karena Race Condition!\n\n",
			expectedBalance-unsafeAcc.Balance)
	} else {
		fmt.Println("⚠️ Kebetulan pas, tapi tetap berbahaya.\n")
	}

	// ----------------------------------------------------------
	// Uji Coba 2: Rekening Aman dengan Mutex
	// ----------------------------------------------------------
	mutexAcc := &MutexAccount{Balance: 0}
	var wg2 sync.WaitGroup

	for i := 0; i < totalTransactions; i++ {
		wg2.Add(1)
		go func() {
			defer wg2.Done()
			mutexAcc.Deposit(depositAmount)
		}()
	}
	wg2.Wait()

	fmt.Println("--- HASIL EKSPERIMEN 2: DENGAN SYNC.MUTEX (AMAN) ---")
	fmt.Printf("Saldo Akhir   : Rp%d\n", mutexAcc.Balance)
	fmt.Printf("✅ SALDO TEPAT 100%% SESUAI EKSPEKTASI: Rp%d\n\n", expectedBalance)

	// ----------------------------------------------------------
	// Uji Coba 3: Rekening Aman dengan Atomic Operations
	// ----------------------------------------------------------
	atomicAcc := &AtomicAccount{Balance: 0}
	var wg3 sync.WaitGroup

	for i := 0; i < totalTransactions; i++ {
		wg3.Add(1)
		go func() {
			defer wg3.Done()
			atomicAcc.Deposit(int64(depositAmount))
		}()
	}
	wg3.Wait()

	fmt.Println("--- HASIL EKSPERIMEN 3: DENGAN SYNC/ATOMIC (AMAN & CEPAT) ---")
	fmt.Printf("Saldo Akhir   : Rp%d\n", atomicAcc.Balance)
	fmt.Printf("✅ SALDO TEPAT 100%% SESUAI EKSPEKTASI: Rp%d\n", expectedBalance)
	fmt.Println("==================================================================")
}
