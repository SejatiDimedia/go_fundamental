package main

import "fmt"

func main() {
	fmt.Println("=== EKSPERIMEN 1: 3-INDEX SLICING ===")
	base := []int{10, 20, 30, 40, 50}

	// a. Slicing biasa: mewarisi sisa kapasitas sampai ujung
	sub1 := base[1:3] // len: 2, cap: 4
	fmt.Printf("sub1: %v (len: %d, cap: %d)\n", sub1, len(sub1), cap(sub1))

	// b. 3-index slicing: mengunci kapasitas agar tidak bocor
	sub2 := base[1:3:3] // len: 2, cap: 2
	fmt.Printf("sub2: %v (len: %d, cap: %d)\n\n", sub2, len(sub2), cap(sub2))

	// Bahaya overwrite pada sub1:
	sub1 = append(sub1, 999)
	fmt.Println("Setelah sub1 di-append 999:")
	fmt.Printf("base ikut rusak: %v\n\n", base)

	fmt.Println("=== EKSPERIMEN 2: RAHASIA APPEND DALAM FUNGSI ===")
	nums := make([]int, 2, 5) // [0, 0]
	tambahElemen(nums)

	fmt.Printf("nums standar (len: 2): %v\n", nums)
	// Kita paksa intip slot ke-3 di underlying array:
	fmt.Printf("nums diintip [:3]     : %v (Angka 999 ternyata ADA di memori!)\n", nums[:3])
}

func tambahElemen(s []int) {
	s = append(s, 999)
}
