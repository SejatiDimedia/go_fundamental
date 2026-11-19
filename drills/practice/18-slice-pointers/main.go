package main

import (
	"fmt"
	"unsafe"
)

func modifyElement(s []int) {
	s[0] = 88
}

func appendElement(s []int) {
	fmt.Printf("  [appendElement awal]  addr: %p, len: %d, cap: %d\n", s, len(s), cap(s))
	s = append(s, 99)
	fmt.Printf("  [appendElement akhir] addr: %p, len: %d, cap: %d, content: %v\n", s, len(s), cap(s), s)
}

func appendElementPtr(s *[]int) {
	*s = append(*s, 99)
}

func main() {
	nums := []int{1, 2, 3}

	fmt.Println("--- Kondisi Awal ---")
	fmt.Printf("nums: %v, len: %d, cap: %d, underlying array ptr: %p\n\n", nums, len(nums), cap(nums), nums)

	fmt.Println("--- 1. Panggil modifyElement(nums) ---")
	modifyElement(nums)
	fmt.Println("1. nums setelah modifyElement:", nums)
	fmt.Println()

	fmt.Println("--- 2. Panggil appendElement(nums) ---")
	appendElement(nums)
	fmt.Println("2. nums setelah appendElement:", nums)
	fmt.Printf("   (Perhatikan: panjang nums di main TETAP %d, isinya TETAP %v)\n\n", len(nums), nums)

	fmt.Println("--- 3. Panggil appendElementPtr(&nums) dengan Pointer ---")
	appendElementPtr(&nums)
	fmt.Println("3. nums setelah appendElementPtr:", nums)
	fmt.Printf("   (Berhasil! len menjadi: %d, isinya: %v)\n", len(nums), nums)

	// Bukti ukuran Slice Header adalah 24 bytes (3 x 8 byte) di arsitektur 64-bit
	fmt.Println("\n--- Struktur Slice Header di RAM ---")
	fmt.Printf("Ukuran Slice Header nums: %d bytes (Pointer Data 8B + Len 8B + Cap 8B)\n", unsafe.Sizeof(nums))
}
