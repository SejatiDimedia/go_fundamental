package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Rule mendefinisikan aturan kelipatan dan kata penggantinya
type Rule struct {
	Divisor int
	Word    string
}

// GenerateFizzBuzz memproses rentang angka dari 'start' sampai 'end' berdasarkan daftar 'rules'
func GenerateFizzBuzz(start, end int, rules []Rule) []string {
	results := make([]string, 0, end-start+1)

	for num := start; num <= end; num++ {
		var builder strings.Builder

		// Cek setiap rule yang aktif
		for _, rule := range rules {
			if num%rule.Divisor == 0 {
				builder.WriteString(rule.Word)
			}
		}

		// Jika tidak ada rule yang cocok (builder masih kosong), gunakan angkanya
		if builder.Len() == 0 {
			results = append(results, strconv.Itoa(num))
		} else {
			results = append(results, builder.String())
		}
	}

	return results
}

func readInput(scanner *bufio.Scanner, label string) string {
	fmt.Print(label)
	scanner.Scan()
	return strings.TrimSpace(scanner.Text())
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("========================================")
	fmt.Println("       🎯 CUSTOM DYNAMIC FIZZBUZZ       ")
	fmt.Println("========================================")

	for {
		fmt.Println("\nPilihan Mode:")
		fmt.Println("1. Mode Standar (Kelipatan 3 = 'Fizz', 5 = 'Buzz')")
		fmt.Println("2. Mode Kustom (Tentukan sendiri kelipatan & katanya)")
		fmt.Println("3. Keluar")

		choice := readInput(scanner, "Pilih mode (1-3): ")
		if choice == "3" || strings.ToLower(choice) == "exit" {
			fmt.Println("Terima kasih! Sampai jumpa.")
			break
		}

		// Input rentang angka
		startStr := readInput(scanner, "Mulai dari angka (start): ")
		endStr := readInput(scanner, "Sampai angka (end): ")

		start, err1 := strconv.Atoi(startStr)
		end, err2 := strconv.Atoi(endStr)

		if err1 != nil || err2 != nil || start > end {
			fmt.Println("❌ Error: Rentang angka tidak valid (pastikan start <= end)!")
			continue
		}

		var activeRules []Rule

		if choice == "1" {
			// Aturan Klasik
			activeRules = []Rule{
				{Divisor: 3, Word: "Fizz"},
				{Divisor: 5, Word: "Buzz"},
			}
		} else if choice == "2" {
			// Aturan Kustom dari Pengguna
			fmt.Println("\n--- Tambah Aturan Kustom ---")
			fmt.Println("Contoh: Kelipatan 3 -> Fizz, Kelipatan 5 -> Buzz, Kelipatan 7 -> Bazz")
			for {
				divStr := readInput(scanner, "Masukkan angka pembagi/kelipatan (atau 'done' jika selesai): ")
				if strings.ToLower(divStr) == "done" {
					break
				}

				divisor, err := strconv.Atoi(divStr)
				if err != nil || divisor <= 0 {
					fmt.Println("❌ Angka pembagi harus berupa integer positif!")
					continue
				}

				word := readInput(scanner, fmt.Sprintf("Kata pengganti untuk kelipatan %d: ", divisor))
				if word == "" {
					fmt.Println("❌ Kata pengganti tidak boleh kosong!")
					continue
				}

				activeRules = append(activeRules, Rule{Divisor: divisor, Word: word})
				fmt.Printf("✅ Aturan ditambahkan: Kelipatan %d -> %q\n", divisor, word)
			}

			if len(activeRules) == 0 {
				fmt.Println("⚠️ Tidak ada aturan yang dimasukkan, menggunakan aturan default (3 & 5).")
				activeRules = []Rule{
					{Divisor: 3, Word: "Fizz"},
					{Divisor: 5, Word: "Buzz"},
				}
			}
		} else {
			fmt.Println("❌ Pilihan tidak valid!")
			continue
		}

		// Eksekusi Logika FizzBuzz
		fmt.Println("\n================ HASIL ================")
		output := GenerateFizzBuzz(start, end, activeRules)

		for i, res := range output {
			fmt.Printf("%-10s", res)
			// Format output rapi: cetak baris baru setiap 5 item
			if (i+1)%5 == 0 {
				fmt.Println()
			}
		}
		fmt.Println("\n=======================================")
	}
}
