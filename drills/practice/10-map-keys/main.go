package main

import "fmt"

// -------------------------------------------------------------
// 1. CONTOH STRUCT
// -------------------------------------------------------------

// Coordinate: SEMUA field-nya comparable (int dan int).
// Struct ini BOLEH dijadikan key map!
type Coordinate struct {
	X, Y int
}

// UserWithTags: Punya field 'Tags []string' (Slice).
// Struct ini DILARANG dijadikan key map karena slice tidak comparable!
type UserWithTags struct {
	Name string
	Tags []string
}

func main() {
	fmt.Println("==================================================")
	fmt.Println("  1. TIPE DATA YANG BOLEH JADI KEY MAP (LEGAL)    ")
	fmt.Println("==================================================")

	// A. Tipe Primitif
	mapString := map[string]int{"Jakarta": 1, "Surabaya": 2}
	mapInt := map[int]string{200: "OK", 404: "Not Found"}
	mapBool := map[bool]string{true: "Aktif", false: "Blokir"}

	// B. Array Berukuran Tetap (Bukan Slice!)
	// [2]int punya panjang tetap dan bisa dibandingkan dengan '=='
	mapArray := map[[2]int]string{
		{0, 0}: "Titik Nol (Origin)",
		{3, 4}: "Titik Koordinat A",
	}

	// C. Struct (yang seluruh field-nya comparable)
	mapStruct := map[Coordinate]string{
		{X: -6, Y: 106}: "Monas Jakarta",
		{X: -7, Y: 110}: "Tugu Jogja",
	}

	// D. Pointer
	// Yang di-hash dan dibandingkan adalah ALAMAT MEMORINYA (0x...)
	val := 100
	mapPointer := map[*int]string{
		&val: "Pointer ke variabel val",
	}

	// E. Channel
	// Yang dibandingkan adalah referensi channel di heap
	ch := make(chan int)
	mapChan := map[chan int]string{
		ch: "Job Queue Channel",
	}

	fmt.Println("✅ mapString:", mapString["Jakarta"])
	fmt.Println("✅ mapInt   :", mapInt[200])
	fmt.Println("✅ mapBool  :", mapBool[true])
	fmt.Println("✅ mapArray :", mapArray[[2]int{0, 0}])
	fmt.Println("✅ mapStruct:", mapStruct[Coordinate{X: -6, Y: 106}])
	fmt.Println("✅ mapPointer:", mapPointer[&val])
	fmt.Println("✅ mapChan  :", mapChan[ch])

	fmt.Println("\n==================================================")
	fmt.Println("  2. TIPE DATA YANG DILARANG JADI KEY (ILLEGAL)   ")
	fmt.Println("==================================================")
	fmt.Println("Semua baris di bawah ini sengaja di-comment.")
	fmt.Println("Jika dibuka, compiler Go akan langsung MENOLAK (Compile Error):")

	// ❌ DILARANG 1: Slice ([]T)
	// var illegal1 map[[]string]string
	// Error: invalid map key type []string

	// ❌ DILARANG 2: Map lain
	// var illegal2 map[map[string]int]string
	// Error: invalid map key type map[string]int

	// ❌ DILARANG 3: Function
	// var illegal3 map[func()]string
	// Error: invalid map key type func()

	// ❌ DILARANG 4: Struct yang mengandung Slice atau Map
	// var illegal4 map[UserWithTags]string
	// Error: invalid map key type UserWithTags
}
