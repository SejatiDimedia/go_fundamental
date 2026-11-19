package main

import "fmt"

type UpdateProductRequest struct {
	Name     string  `json:"name"`
	Price    *int    `json:"price"`
	Discount *int    `json:"discount"`
	IsActive *bool   `json:"is_active"`
	Bio      *string `json:"bio"`
}

// TRIK SAKTI SENIOR: Generic Pointer Helper
// Mengubah value literal apapun (int, string, bool) menjadi pointer dalam 1 baris!
func Ptr[T any](v T) *T {
	return &v
}

func main() {
	// ==========================================
	// ❌ CARA PEMULA (Melelahkan & Banyak Variabel Temp):
	// ==========================================
	// price := 100000
	// disc := 10
	// active := true
	// bio := "Produk unggulan"
	// reqJunior := UpdateProductRequest{
	// 	Name:     "Kemeja",
	// 	Price:    &price,
	// 	Discount: &disc,
	// 	IsActive: &active,
	// 	Bio:      &bio,
	// }

	// ==========================================
	// ⚡ CARA SENIOR DENGAN TRIK Ptr[T]:
	// ==========================================
	reqSenior := UpdateProductRequest{
		Name:     "Kemeja",
		Price:    Ptr(100000),        // Langsung int jadi *int!
		Discount: Ptr(10),            // Langsung int jadi *int!
		IsActive: Ptr(true),          // Langsung bool jadi *bool!
		Bio:      Ptr("Produk unggulan"), // Langsung string jadi *string!
	}

	fmt.Println("=== Sukses Mengisi Struct Pointer Tanpa Variabel Temp ===")
	fmt.Printf("Name: %s\n", reqSenior.Name)
	fmt.Printf("Price: %d\n", *reqSenior.Price)
	fmt.Printf("Discount: %d%%\n", *reqSenior.Discount)
	fmt.Printf("IsActive: %t\n", *reqSenior.IsActive)
	fmt.Printf("Bio: %s\n", *reqSenior.Bio)
}
