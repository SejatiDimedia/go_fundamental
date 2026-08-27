package main

import (
	"cmp"
	"fmt"
)

// A. Generic Function: Mencari nilai terkecil/minimum dari tipe numerik atau string apa pun
func Min[T cmp.Ordered](<a, b T>) T {
	if a < b {
		return a
	}
	return b
}

// B. Generic In-Memory Repository: Wadah penyimpanan data apa saja (Produk, User, Order)
type Repository[T any] struct {
	storage map[string]T
}

func NewRepository[T any]() *Repository[T] {
	return &Repository[T]{
		storage: make(map[string]T),
	}
}

func (r *Repository[T]) Save(id string, item T) {
	r.storage[id] = item
}

func (r *Repository[T]) FindByID(id string) (T, bool) {
	item, exists := r.storage[id]
	return item, exists
}

// Contoh Struct Domain
type User struct {
	Username string
	Role     string
}

type ProductItem struct {
	Title string
	Price int
}

func main() {
	// 1. Uji Generic Min Function
	fmt.Println("Min Integer :", Min(100, 45))        // Otomatis tipe int -> 45
	fmt.Println("Min Float   :", Min(12.5, 3.8))      // Otomatis tipe float64 -> 3.8
	fmt.Println("Min String  :", Min("Budi", "Andi")) // Mengurutkan alfabet -> "Andi"

	// 2. Uji Generic Repository untuk User
	userRepo := NewRepository[User]()
	userRepo.Save("USR-1", User{Username: "john_doe", Role: "Admin"})

	if u, ok := userRepo.FindByID("USR-1"); ok {
		fmt.Printf("\nUser Ditemukan: %s (%s)\n", u.Username, u.Role)
	}

	// 3. Uji Generic Repository yang SAMA untuk ProductItem
	productRepo := NewRepository[ProductItem]()
	productRepo.Save("PRD-1", ProductItem{Title: "Monitor 4K", Price: 4500000})

	if p, ok := productRepo.FindByID("PRD-1"); ok {
		fmt.Printf("Produk Ditemukan: %s (Rp%d)\n", p.Title, p.Price)
	}
}