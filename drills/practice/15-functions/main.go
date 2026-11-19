package main

import "fmt"

// TODO 1: Variadic Function
// Buat fungsi 'CalculateStats(numbers ...int) (min int, max int, sum int)'
// - Gunakan NAMED RETURN VALUES (min int, max int, sum int)
// - Jika numbers kosong (len == 0), return 0, 0, 0
// - Cari nilai minimum, maximum, dan jumlah total seluruh angka
// - Return secara eksplisit: return min, max, sum
func CalculateStats(numbers ...int) (min int, max int, sum int) {
	if len(numbers) == 0 {
		return 0, 0, 0
	}

	min = numbers[0]
	max = numbers[0]

	for _, number := range numbers {
		if number > max {
			max = number
		}

		if number < min {
			min = number
		}
		sum += number
	}

	return min, max, sum
}

// TODO 2: Closure / Function Factory (Generator ID Transaksi)
// Buat fungsi 'CreateIDGenerator(prefix string) func() string'
//   - Fungsi ini mengembalikan fungsi anonim: func() string
//   - Simpan sebuah variabel lokal 'counter := 0' di dalam CreateIDGenerator
//   - Setiap kali fungsi anonim yang di-return itu dipanggil:
//     counter akan bertambah +1
//     kembalikan string format: fmt.Sprintf("%s-%03d", prefix, counter)
//     Contoh jika prefix "TRX": "TRX-001", "TRX-002", dst.
func CreateIDGenerator(prefix string) func() string {
	counter := 0
	return func() string {
		counter++
		return fmt.Sprintf("%s-%03d", prefix, counter)
	}
}

func main() {
	fmt.Println("=== 1. UJI VARIADIC FUNCTION & SLICE UNPACKING ===")
	scores := []int{85, 92, 78, 95, 88}

	// Panggil CalculateStats dengan membongkar slice scores (scores...)
	min, max, sum := CalculateStats(scores...)
	fmt.Printf("Min: %d, Max: %d, Sum: %d\n", min, max, sum)
	// Ekspektasi: Min: 78, Max: 95, Sum: 438

	fmt.Println("\n=== 2. UJI CLOSURE / FUNCTION GENERATOR ===")
	// Buat dua generator independen:
	trxGen := CreateIDGenerator("TRX")
	usrGen := CreateIDGenerator("USR")
	//
	fmt.Println(trxGen()) // Ekspektasi: TRX-001
	fmt.Println(trxGen()) // Ekspektasi: TRX-002
	fmt.Println(usrGen()) // Ekspektasi: USR-001
	fmt.Println(trxGen()) // Ekspektasi: TRX-003
	// (Perhatikan: trxGen dan usrGen memiliki state counter masing-masing yang terisolasi!)
}
