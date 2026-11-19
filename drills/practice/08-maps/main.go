package main

import "fmt"

// TODO 1: Buat fungsi AddStock(inventory map[string]int, item string, qty int)
// - Tambahkan qty ke inventory[item] (gunakan += qty)
// func AddStock ...
func AddStock(inventory map[string]int, item string, qty int) {
	inventory[item] += qty
}

// TODO 2: Buat fungsi GetStock(inventory map[string]int, item string) (int, bool)
// - Gunakan comma-ok idiom: val, ok := inventory[item]
// - Return val, ok
// func GetStock ...
func GetStock(inventory map[string]int, item string) (int, bool) {
	val, ok := inventory[item]
	return val, ok
}

// TODO 3: Buat fungsi DeleteStock(inventory map[string]int, item string)
// - Hapus item menggunakan: delete(inventory, item)
// func DeleteStock ...
func DeleteStock(inventory map[string]int, item string) {
	delete(inventory, item)
}

func main() {
	// TODO 4:
	// a. Inisialisasi map stok dengan make(map[string]int)
	var inventory = make(map[string]int)
	// b. Panggil AddStock(stok, "apel", 10)
	AddStock(inventory, "Apel", 10)
	// c. Panggil AddStock(stok, "apel", 5)  // akumulasi jadi 15
	AddStock(inventory, "Apel", 5)
	// d. Uji GetStock:
	//    - Cek "apel"   -> ekspektasi: 15 true
	//    - Cek "durian" -> ekspektasi: 0 false
	fmt.Println(GetStock(inventory, "Apel"))
	fmt.Println(GetStock(inventory, "Durian"))
	// e. Hapus "apel" dengan DeleteStock
	DeleteStock(inventory, "Apel")
	// f. Cek lagi "apel" -> ekspektasi: 0 false
	fmt.Println(GetStock(inventory, "Apel"))
}
