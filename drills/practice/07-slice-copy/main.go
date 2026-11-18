package main

import "fmt"

// TODO 1: Buat fungsi CloneAndAppend(original []string, item string) []string
// - Buat slice baru: make([]string, len(original))
// - Salin data: copy(baru, original)
// - Append item: baru = append(baru, item)
// - Return baru
// func CloneAndAppend ...
func CloneAndAppend(original []string, item string) []string {
	var baru = make([]string, len(original))
	copy(baru, original)
	baru = append(baru, item)
	return baru
}

// TODO 2: Buat fungsi CutSafe(original []int, low, high int) []int
// - Potong dengan 3-index slicing: original[low:high:high]
// - Return hasil potongannya
// func CutSafe ...
func CutSafe(original []int, low, high int) []int {
	var baru = original[low:high:high]
	return baru
}

func main() {
	// TODO 3:
	// a. Uji CloneAndAppend:
	//    users := []string{"Budi", "Siti"}
	//    newUsers := CloneAndAppend(users, "Andi")
	//    Cetak users    -> harus tetap: [Budi Siti]
	//    Cetak newUsers -> harus: [Budi Siti Andi]
	var users = []string{"Budi", "Siti"}
	var newUsers = CloneAndAppend(users, "Andi")
	fmt.Println(users)
	fmt.Println(newUsers)

	// b. Uji CutSafe:
	//    nums := []int{10, 20, 30, 40, 50}
	//    sub := CutSafe(nums, 1, 3) // ambil index 1 dan 2 (angka 20, 30)
	//    Cetak sub -> [20 30]
	//    Cetak len(sub) dan cap(sub) -> harus len: 2, cap: 2!
	var nums = []int{10, 20, 30, 40, 50}
	var sub = CutSafe(nums, 1, 3)
	fmt.Println(sub)
	fmt.Println(len(sub))
	fmt.Println(cap(sub))

	fmt.Println("Silakan lengkapi latihan di atas!")
}
