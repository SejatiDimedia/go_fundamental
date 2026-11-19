package main

import "fmt"

func modifyValue(p *int) {
	fmt.Printf("[modifyValue] p menunjuk ke alamat: %p\n", p)
	*p = 99 // Mengubah isi di alamat p
}

func modifyPointer(p *int) {
	fmt.Printf("[modifyPointer awal] p menunjuk ke alamat: %p\n", p)
	newVal := 100
	p = &newVal // Mengubah ke mana p menunjuk (kertas alamat lokal)
	fmt.Printf("[modifyPointer akhir] p sekarang menunjuk ke alamat baru: %p (newVal: %d)\n", p, *p)
}

func main() {
	x := 10
	fmt.Printf("Alamat asli x di main: %p, Nilai awal x: %d\n\n", &x, x)

	fmt.Println("--- 1. Panggil modifyValue(&x) ---")
	modifyValue(&x)
	fmt.Printf("Nilai x setelah modifyValue: %d (Alamat x: %p)\n\n", x, &x)

	fmt.Println("--- 2. Panggil modifyPointer(&x) ---")
	modifyPointer(&x)
	fmt.Printf("Nilai x setelah modifyPointer: %d (Alamat x: %p)\n", x, &x)
}
