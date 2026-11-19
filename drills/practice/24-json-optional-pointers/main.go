package main

import (
	"encoding/json"
	"fmt"
)

// Request A: Menggunakan int biasa (Zero Value Trap)
type UpdateProductValue struct {
	Name  string `json:"name"`
	Price int    `json:"price"` // Zero value: 0
}

// Request B: Menggunakan *int (Optional Pointer Pattern)
type UpdateProductPtr struct {
	Name  string `json:"name"`
	Price *int   `json:"price"` // Zero value: nil
}

func main() {
	// JSON 1: Klien hanya ingin update nama, TIDAK kirim field price
	json1 := `{"name": "Kemeja Putih"}`

	// JSON 2: Klien sengaja pasang diskon 100% (GRATIS, price = 0)
	json2 := `{"name": "Promo Gratis", "price": 0}`

	fmt.Println("=== 1. MASALAH: Menggunakan int Biasa ===")
	var reqVal1, reqVal2 UpdateProductValue
	json.Unmarshal([]byte(json1), &reqVal1)
	json.Unmarshal([]byte(json2), &reqVal2)

	fmt.Printf("JSON 1 -> Price: %d (Apakah klien pasang harga 0 atau tidak kirim? TIDAK TAHU!)\n", reqVal1.Price)
	fmt.Printf("JSON 2 -> Price: %d\n\n", reqVal2.Price)

	fmt.Println("=== 2. SOLUSI: Menggunakan *int (Pointer Pattern) ===")
	var reqPtr1, reqPtr2 UpdateProductPtr
	json.Unmarshal([]byte(json1), &reqPtr1)
	json.Unmarshal([]byte(json2), &reqPtr2)

	if reqPtr1.Price == nil {
		fmt.Printf("JSON 1 -> Price: nil => Klien TIDAK MENGIRIM field price (jangan ubah di DB!)\n")
	}

	if reqPtr2.Price != nil {
		fmt.Printf("JSON 2 -> Price: %d  => Klien EKSPLISIT mengubah harga jadi 0!\n", *reqPtr2.Price)
	}
}
