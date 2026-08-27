package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func calculate(num1 float64, operator string, num2 float64) (float64, error) {
	switch operator {
	case "+":
		return num1 + num2, nil
	case "-":
		return num1 - num2, nil
	case "*":
		return num1 * num2, nil
	case "/":
		if num2 == 0 {
			return 0, errors.New("Tidak bisa membagi dengan angka 0")
		}
		return num1 / num2, nil
	default:
		return 0, errors.New("Operator tidak dikenal")
	}
}

func promptInput(scanner *bufio.Scanner, label string) string {
	fmt.Print(label)
	scanner.Scan()
	return strings.TrimSpace(scanner.Text())
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("========================================")
	fmt.Println("       🧮 GOLANG CLI CALCULATOR        ")
	fmt.Println("========================================")

	for {
		input1 := promptInput(scanner, "\nMasukkan angkat pertama (atau ketik 'exit' untuk keluar): ")
		if strings.ToLower(input1) == "exit" {
			fmt.Println("Terima kasih! Sampai jumpa.")
			break
		}

		num1, err := strconv.ParseFloat(input1, 64)
		if err != nil {
			fmt.Println("❌ Error: Input angka pertama harus berupa angka yang valid")
			continue
		}

		operator := promptInput(scanner, "Pilih operator (+, -, *, /): ")

		input2 := promptInput(scanner, "Masukkan angka kedua: ")
		num2, err := strconv.ParseFloat(input2, 64)
		if err != nil {
			fmt.Println("❌ Error: Input angka kedua harus berupa angka yang valid")
			continue
		}

		result, err := calculate(num1, operator, num2)
		if err != nil {
			fmt.Println("❌ Gagal menghitung: %s\n", err.Error())
		} else {
			fmt.Println("-------------------------------------")
			fmt.Printf("✅ Hasil: %.2f %s %.2f = %.2f\n", num1, operator, num2, result)
			fmt.Println("----------------------------------------")
		}
	}

}
