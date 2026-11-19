package main

import (
	"fmt"
	"runtime"
)

type Account struct {
	ID      string
	Balance int
}

func printAlloc() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("Allocated Memory: %d MB\n", m.Alloc/1024/1024)
}

func main() {
	fmt.Println("=== BAGIAN 1: STRUCT IN MAP MUTABILITY ===")
	// TODO 1:
	// Buat map bernama 'accounts' dengan tipe map[string]*Account (value-nya pointer!)
	// Masukkan akun "ACC-01" dengan Balance: 50000
	// Tambahkan Balance-nya menjadi 75000 secara langsung: accounts["ACC-01"].Balance += 25000
	// Print Balance akhir dari "ACC-01"
	accounts := map[string]*Account{
		"ACC-01": &Account{
			ID:      "ACC-01",
			Balance: 50000,
		},
	}
	accounts["ACC-01"].Balance += 25000
	fmt.Println(accounts["ACC-01"].Balance)

	// TODO 2:
	// Sekarang gunakan map bernilai struct biasa (BUKAN pointer):
	// accountsVal := map[string]Account{"ACC-02": {ID: "ACC-02", Balance: 10000}}
	// Update Balance ACC-02 menjadi 20000 dengan teknik Copy -> Mutate -> Re-assign!
	// Print Balance akhir dari "ACC-02"
	accountsVal := map[string]Account{
		"ACC-02": {
			ID:      "ACC-02",
			Balance: 10000,
		},
	}
	a := accountsVal["ACC-02"]
	a.Balance = 20000
	accountsVal["ACC-02"] = a
	fmt.Println(accountsVal["ACC-02"].Balance)

	fmt.Println("\n=== BAGIAN 2: MAP BUCKET MEMORY RETENTION ===")
	// Kita simulasikan alokasi memori map besar
	fmt.Print("1. Kondisi Awal: ")
	printAlloc()

	largeMap := make(map[int][128]byte)
	for i := 0; i < 500_000; i++ {
		largeMap[i] = [128]byte{}
	}
	runtime.GC()
	fmt.Print("2. Setelah diisi 500.000 elemen: ")
	printAlloc()

	// Hapus semua elemen
	for i := 0; i < 500_000; i++ {
		delete(largeMap, i)
	}
	runtime.GC()
	fmt.Printf("3. Setelah semua elemen di-DELETE (len=%d): ", len(largeMap))
	printAlloc()
	runtime.KeepAlive(largeMap) // Menjaga largeMap tetap hidup saat pengecekan memori

	// TODO 3:
	// Buat largeMap mengembalikan memorinya ke GC!
	// (Petunjuk: Re-assign largeMap dengan make baru atau nil)
	// largeMap = ...
	largeMap = make(map[int][128]byte)

	runtime.GC()
	fmt.Printf("4. Setelah di-reassign map baru (len=%d): ", len(largeMap))
	printAlloc()
	runtime.KeepAlive(largeMap)
}
