package main

import "fmt"

// TODO 1: Buat fungsi RemoveOrdered(s []int, index int) []int
// - Validasi: jika index < 0 atau index >= len(s), return s langsung
// - Hapus elemen dengan: append(s[:index], s[index+1:]...)
// - Return hasil slice
// func RemoveOrdered ...
func RemoveOrdered(s []int, index int) []int {
	return append(s[:index], s[index+1:]...)
}

// TODO 2: Buat fungsi RemoveFast(s []int, index int) []int (O(1) Unordered Delete)
// - Validasi: jika index < 0 atau index >= len(s), return s langsung
// - Timpa s[index] dengan elemen terakhir s[len(s)-1]
// - Potong elemen terakhir: s = s[:len(s)-1]
// - Return s
// func RemoveFast ...
func RemoveFast(s []int, index int) []int {
	if index < 0 || index >= len(s) {
		return s
	}
	s[index] = s[len(s)-1] // s[1] = s[5-1]
	return s[:len(s)-1]
}

func main() {
	// TODO 3:
	// a. Uji RemoveOrdered:
	//    list1 := []int{10, 20, 30, 40, 50}
	//    list1 = RemoveOrdered(list1, 2) // Hapus index 2 (angka 30)
	//    Cetak list1 -> ekspektasi: [10 20 40 50]
	list1 := []int{10, 20, 30, 40, 50}
	list1 = RemoveOrdered(list1, 2)
	fmt.Println(list1)

	// b. Uji RemoveFast:
	//    list2 := []int{10, 20, 30, 40, 50}
	//    list2 = RemoveFast(list2, 1) // Hapus index 1 (angka 20, ditimpa 50)
	//    Cetak list2 -> ekspektasi: [10 50 30 40]
	list2 := []int{10, 20, 30, 40, 50}
	list2 = RemoveFast(list2, 1)
	fmt.Println(list2)

}
