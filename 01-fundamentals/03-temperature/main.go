package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Custom type untuk memperjelas makna domain data
type Temperature float64
type Scale string

const (
	Celsius    Scale = "C"
	Fahrenheit Scale = "F"
	Kelvin     Scale = "K"
	Reamur     Scale = "R"
)

// toCelsius mengubah suhu dari skala apa pun ke skala basis (Celsius)
func toCelsius(val Temperature, from Scale) (Temperature, error) {
	switch from {
	case Celsius:
		return val, nil
	case Fahrenheit:
		return Temperature((float64(val) - 32) * 5 / 9), nil
	case Kelvin:
		if val < 0 {
			return 0, errors.New("suhu Kelvin tidak boleh kurang dari 0 (Absolute Zero)")
		}
		return Temperature(float64(val) - 273.15), nil
	case Reamur:
		return Temperature(float64(val) * 5 / 4), nil
	default:
		return 0, fmt.Errorf("skala asal %q tidak didukung", from)
	}
}

// fromCelsius mengubah suhu dari basis Celsius ke skala target
func fromCelsius(celsiusVal Temperature, to Scale) (Temperature, error) {
	switch to {
	case Celsius:
		return celsiusVal, nil
	case Fahrenheit:
		return Temperature((float64(celsiusVal) * 9 / 5) + 32), nil
	case Kelvin:
		kVal := float64(celsiusVal) + 273.15
		if kVal < 0 {
			return 0, errors.New("hasil konversi menghasilkan nilai di bawah Absolute Zero (0 Kelvin)")
		}
		return Temperature(kVal), nil
	case Reamur:
		return Temperature(float64(celsiusVal) * 4 / 5), nil
	default:
		return 0, fmt.Errorf("skala tujuan %q tidak didukung", to)
	}
}

// convertTemperature adalah fungsi utama penghubung (pipeline konversi)
func convertTemperature(val Temperature, from Scale, to Scale) (Temperature, error) {
	celsius, err := toCelsius(val, from)
	if err != nil {
		return 0, err
	}
	return fromCelsius(celsius, to)
}

func readLine(scanner *bufio.Scanner, prompt string) string {
	fmt.Print(prompt)
	scanner.Scan()
	return strings.TrimSpace(scanner.Text())
}

func parseScale(input string) (Scale, error) {
	upper := strings.ToUpper(input)
	switch Scale(upper) {
	case Celsius, Fahrenheit, Kelvin, Reamur:
		return Scale(upper), nil
	default:
		return "", fmt.Errorf("skala %q tidak valid. Pilihan: C, F, K, R", input)
	}
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("========================================")
	fmt.Println("     🌡️  GOLANG TEMPERATURE CONVERTER   ")
	fmt.Println("   Skala yang didukung: C, F, K, R      ")
	fmt.Println("========================================")

	for {
		// 1. Input Nilai Suhu
		inputVal := readLine(scanner, "\nMasukkan nilai suhu (atau 'exit' untuk selesai): ")
		if strings.ToLower(inputVal) == "exit" {
			fmt.Println("Terima kasih! Sampai jumpa.")
			break
		}

		rawFloat, err := strconv.ParseFloat(inputVal, 64)
		if err != nil {
			fmt.Println("❌ Error: Masukkan angka suhu yang valid!")
			continue
		}
		tempVal := Temperature(rawFloat)

		// 2. Input Skala Asal
		fromInput := readLine(scanner, "Dari skala apa? (C/F/K/R): ")
		fromScale, err := parseScale(fromInput)
		if err != nil {
			fmt.Printf("❌ Error: %s\n", err.Error())
			continue
		}

		// 3. Input Skala Target
		toInput := readLine(scanner, "Konversi ke skala apa? (C/F/K/R): ")
		toScale, err := parseScale(toInput)
		if err != nil {
			fmt.Printf("❌ Error: %s\n", err.Error())
			continue
		}

		// 4. Eksekusi Konversi
		result, err := convertTemperature(tempVal, fromScale, toScale)
		if err != nil {
			fmt.Printf("❌ Gagal konversi: %s\n", err.Error())
		} else {
			fmt.Println("----------------------------------------")
			fmt.Printf("✅ Hasil: %.2f °%s = %.2f °%s\n", tempVal, fromScale, result, toScale)
			fmt.Println("----------------------------------------")
		}
	}
}
