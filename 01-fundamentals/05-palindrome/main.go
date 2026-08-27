package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"unicode"
)

// cleanString membersihkan teks dari tanda baca & spasi, serta mengubah ke lowercase rune
func cleanString(input string) []rune {
	var cleaned []rune
	for _, r := range input {
		// Hanya ambil karakter huruf atau angka
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			cleaned = append(cleaned, unicode.ToLower(r))
		}
	}
	return cleaned
}

// IsPalindrome mengecek apakah deretan rune simetris menggunakan teknik Two Pointers
func IsPalindrome(input string) (bool, string) {
	runes := cleanString(input)

	if len(runes) == 0 {
		return false, "Teks kosong atau tidak mengandung huruf/angka"
	}

	left := 0
	right := len(runes) - 1

	for left < right {
		if runes[left] != runes[right] {
			return false, fmt.Sprintf("Karakter '%c' pada posisi awal tidak cocok dengan '%c' pada posisi akhir", runes[left], runes[right])
		}
		left++
		right--
	}

	return true, "Semua karakter simetris dari kiri dan kanan"
}

// ReverseString membalikkan teks dengan aman berbasis rune (mendukung karakter Unicode/Aksen)
func ReverseString(input string) string {
	runes := []rune(input)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

func readInput(scanner *bufio.Scanner, label string) string {
	fmt.Print(label)
	scanner.Scan()
	return strings.TrimSpace(scanner.Text())
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("========================================")
	fmt.Println("     🔄 UNICODE PALINDROME CHECKER      ")
	fmt.Println("========================================")

	for {
		input := readInput(scanner, "\nMasukkan kata/kalimat (atau 'exit' untuk selesai): ")
		if strings.ToLower(input) == "exit" {
			fmt.Println("Terima kasih! Sampai jumpa.")
			break
		}

		if input == "" {
			fmt.Println("⚠️ Teks tidak boleh kosong!")
			continue
		}

		isPal, reason := IsPalindrome(input)
		reversed := ReverseString(input)

		fmt.Println("----------------------------------------")
		fmt.Printf("Teks Asli   : %s\n", input)
		fmt.Printf("Teks Dibalik: %s\n", reversed)
		if isPal {
			fmt.Println("Status      : ✅ PALINDROME!")
		} else {
			fmt.Println("Status      : ❌ BUKAN PALINDROME")
		}
		fmt.Printf("Penjelasan  : %s\n", reason)
		fmt.Println("----------------------------------------")
	}
}
