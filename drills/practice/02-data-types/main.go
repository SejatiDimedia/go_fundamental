package main

import "fmt"

// TODO 1: Buat fungsi ChangeInitial(text string, newChar rune) string
// - Tangani jika text kosong ("") -> return text langsung
// - Konversi text ke []rune
// - Ubah index 0 dengan newChar
// - Kembalikan dalam bentuk string
// func ChangeInitial ...
func ChangeInitial(text string, newChar rune) string {
	if len(text) == 0 {
		return text
	}
	var runes = []rune(text)
	runes[0] = newChar
	return string(runes)
}

// TODO 2: Buat fungsi CountRunes(text string) int
// - Kembalikan jumlah karakter Unicode asli (bisa pakai len([]rune(text)))
// func CountRunes ...
func CountRunes(text string) int {
	return len([]rune(text))
}

func main() {
	// TODO 3:
	// a. Panggil ChangeInitial("Golang", 'J') dan cetak hasilnya (ekspektasi: "Jolang")
	fmt.Println(ChangeInitial("Golang", 'J'))
	// b. Uji string "Go™":
	//    - Cetak jumlah byte-nya menggunakan len(...) -> ekspektasi: 5
	//    - Cetak jumlah karakter aslinya menggunakan CountRunes(...) -> ekspektasi: 3
	fmt.Println(len("Go™"))
	fmt.Println(CountRunes("Go™"))

	fmt.Println("Silakan lengkapi kode di atas!")
}
