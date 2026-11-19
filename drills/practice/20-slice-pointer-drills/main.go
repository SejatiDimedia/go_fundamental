package main

import "fmt"

// Percobaan 1: Menggunakan Slice biasa
func ClearCartA(cart []string) {
	cart = cart[:0] // Mengubah panjang slice menjadi 0
}

// Percobaan 2: Menggunakan Pointer to Slice
func ClearCartB(cart *[]string) {
	*cart = (*cart)[:0] // Mengubah panjang slice asli menjadi 0
}

func main() {
	cart1 := []string{"Buku", "Laptop", "Mouse"}
	fmt.Println("--- Kasus 1: Panggil ClearCartA (tanpa pointer) ---")
	fmt.Printf("Sebelum: %v (len: %d)\n", cart1, len(cart1))
	ClearCartA(cart1)
	fmt.Printf("Sesudah: %v (len: %d)\n\n", cart1, len(cart1))

	cart2 := []string{"Buku", "Laptop", "Mouse"}
	fmt.Println("--- Kasus 2: Panggil ClearCartB (dengan pointer) ---")
	fmt.Printf("Sebelum: %v (len: %d)\n", cart2, len(cart2))
	ClearCartB(&cart2)
	fmt.Printf("Sesudah: %v (len: %d)\n", cart2, len(cart2))
}
