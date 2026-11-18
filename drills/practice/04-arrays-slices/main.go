package main

import "fmt"

// TODO 1: Buat fungsi FilterEven(numbers []int) []int
// - Buat slice penampung dengan kapasitas awal: make([]int, 0, len(numbers))
// - Loop numbers, jika n % 2 == 0, append ke penampung
// - Return penampung
// func FilterEven ...
func FilterEven(numbers []int) []int {
	awal := make([]int, 0, len(numbers))
	for _, n := range numbers {
		if n%2 == 0 {
			awal = append(awal, n)
		}
	}
	return awal
}

// TODO 2: Buat fungsi SafeCopyArray(arr [5]int) []int
// - Buat slice baru dengan panjang dan kapasitas 5: make([]int, len(arr))
// - Copy data dari arr[:] ke slice baru tersebut menggunakan copy(...)
// - Return slice baru
// func SafeCopyArray ...
func SafeCopyArray(arr [5]int) []int {
	newArr := make([]int, len(arr))
	copy(newArr, arr[:])
	return newArr
}

func main() {
	// TODO 3:
	// a. Uji FilterEven([]int{1, 2, 3, 4, 5, 6})
	//    Cetak hasilnya -> ekspektasi: [2 4 6]
	fmt.Println(FilterEven([]int{1, 2, 3, 4, 5, 6}))

	// b. Uji SafeCopyArray:
	//    origArr := [5]int{10, 20, 30, 40, 50}
	//    copiedSlice := SafeCopyArray(origArr)
	//    copiedSlice[0] = 999
	//    Cetak origArr (harus tetap [10 20 30 40 50])
	//    Cetak copiedSlice (harus [999 20 30 40 50])
	origArr := [5]int{10, 20, 30, 40, 50}
	copiedSlice := SafeCopyArray(origArr)
	copiedSlice[0] = 999
	fmt.Println(origArr)
	fmt.Println(copiedSlice)

	fmt.Println("Silakan selesaikan latihan di atas!")
}
