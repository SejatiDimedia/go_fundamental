package main

import "fmt"

// Percobaan 1: Gagal mengubah pointer pemanggil
func allocateFail(p *int) {
	val := 42
	p = &val // Hanya mengubah kertas alamat lokal p
}

// Percobaan 2: Berhasil menggunakan Pointer to Pointer (**int)
func allocateSuccess(pp **int) {
	val := 42
	*pp = &val // Mendatangi variabel pointer di main dan mengisinya dengan alamat val
}

func main() {
	var ptr *int // Nilai awal: nil (belum menunjuk ke mana-mana)
	fmt.Println("--- Kondisi Awal ---")
	fmt.Printf("ptr: %v\n\n", ptr)

	fmt.Println("--- 1. Panggil allocateFail(ptr) ---")
	allocateFail(ptr)
	fmt.Printf("ptr setelah allocateFail: %v (Apakah berubah? TIDAK, tetap nil!)\n\n", ptr)

	fmt.Println("--- 2. Panggil allocateSuccess(&ptr) ---")
	allocateSuccess(&ptr)
	fmt.Printf("ptr setelah allocateSuccess: %v\n", ptr)
	fmt.Printf("Nilai yang ditunjuk (*ptr): %d\n", *ptr)
}
