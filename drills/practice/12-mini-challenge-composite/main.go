package main

import (
	"encoding/json"
	"fmt"
)

// ============================================================================
// SPESIFIKASI DATA MODEL
// ============================================================================

// 1. Struct Item
// - SKU: string, tag json: "sku"
// - Name: string, tag json: "name"
// - Price: int, tag json: "price"
type Item struct {
	SKU   string `json:"sku"`
	Name  string `json:"name"`
	Price int    `json:"price"`
}

// 2. Struct CartItem (Wajib gunakan Struct Embedding!)
// - Embed struct Item (Anonymous field)
// - Qty: int, tag json: "qty"
type CartItem struct {
	// TODO 1: Embed struct Item di sini
	// TODO 2: Tambahkan field Qty
	Item
	Qty int `json:"qty"`
}

// 3. Struct Cart
// - UserID: string, tag json: "user_id"
// - Items: []CartItem, tag json: "items"
// - InternalPromoCode: string, tag json: "-" (RAHASIA, tidak boleh keluar di JSON!)
type Cart struct {
	UserID            string     `json:"user_id"`
	Items             []CartItem `json:"items"`
	InternalPromoCode string     `json:"-"`
}

// ============================================================================
// LOGIKA BISNIS
// ============================================================================

// TODO 3: Implementasikan AddToCart
// Aturan:
//  1. Cek apakah sku ada di dalam catalog (gunakan comma-ok idiom pada map!).
//     Jika tidak ada di catalog, return false.
//  2. Jika ada, periksa apakah barang tersebut sudah ada di dalam cart.Items:
//     - Jika sudah ada, cukup tambahkan qty-nya (cart.Items[i].Qty += qty).
//     - Jika belum ada, buat CartItem baru lalu append ke cart.Items.
//  3. Return true jika berhasil.
func AddToCart(catalog map[string]Item, cart *Cart, sku string, qty int) bool {
	existingItem, ok := catalog[sku]
	if !ok {
		return false
	}

	for i := range cart.Items {
		if cart.Items[i].SKU == sku {
			cart.Items[i].Qty += qty
			return true
		}
	}

	cart.Items = append(cart.Items, CartItem{
		Item: existingItem,
		Qty:  qty,
	})
	return true
}

// TODO 4: Implementasikan RemoveFromCart
// Aturan:
// 1. Cari item dengan SKU yang sesuai di cart.Items.
// 2. Hapus elemen tersebut dari slice cart.Items (gunakan teknik delete slice yang sudah kamu pelajari).
// 3. Return true jika berhasil ditemukan & dihapus, false jika tidak ada.
func RemoveFromCart(cart *Cart, sku string) bool {
	for i, item := range cart.Items {
		if item.SKU == sku {
			cart.Items = append(cart.Items[:i], cart.Items[i+1:]...) // cart.Items[0:1] + cart.Items[2:4]
			return true
		}
	}
	return false
}

// TODO 5: Implementasikan CalculateTotal
// Hitung total belanja: jumlahkan (Price * Qty) dari semua item di keranjang.
func CalculateTotal(cart Cart) int {
	total := 0
	for _, item := range cart.Items {
		total += item.Price * item.Qty
	}
	return total
}

func main() {
	fmt.Println("=== MINI CHALLENGE TOPIK 2: SHOPPING CART ===")

	// 1. Siapkan Katalog Barang (Map)
	catalog := map[string]Item{
		"SKU-01": {SKU: "SKU-01", Name: "Kopi Arabika", Price: 30000},
		"SKU-02": {SKU: "SKU-02", Name: "Susu Segar", Price: 20000},
		"SKU-03": {SKU: "SKU-03", Name: "Cokelat Bubuk", Price: 25000},
	}

	// 2. Inisialisasi Keranjang Belanja
	cart := Cart{
		UserID:            "USR-88",
		Items:             make([]CartItem, 0, 5), // Pre-allocate kapasitas 5
		InternalPromoCode: "DISKON-INTERNAL-RAHASIA",
	}

	// 3. Uji AddToCart
	fmt.Println("1. Tambah SKU-01 (Qty: 2):", AddToCart(catalog, &cart, "SKU-01", 2)) // Harusnya true
	fmt.Println("2. Tambah SKU-02 (Qty: 1):", AddToCart(catalog, &cart, "SKU-02", 1)) // Harusnya true
	fmt.Println("3. Tambah SKU-01 (Qty: 3):", AddToCart(catalog, &cart, "SKU-01", 3)) // Harusnya true (Qty SKU-01 jadi 5)
	fmt.Println("4. Tambah SKU-99 (BODONG):", AddToCart(catalog, &cart, "SKU-99", 1)) // Harusnya false

	// 4. Uji CalculateTotal
	// Perhitungan: (SKU-01: 30.000 * 5) + (SKU-02: 20.000 * 1) = 150.000 + 20.000 = 170.000
	total := CalculateTotal(cart)
	fmt.Printf("5. Total Belanja: Rp %d (Ekspektasi: Rp 170000)\n", total)

	// 5. Uji RemoveFromCart
	fmt.Println("6. Hapus SKU-02:", RemoveFromCart(&cart, "SKU-02")) // Harusnya true
	totalSetelahHapus := CalculateTotal(cart)
	// Total sekarang: SKU-01 saja (30.000 * 5 = 150.000)
	fmt.Printf("7. Total Setelah Hapus SKU-02: Rp %d (Ekspektasi: Rp 150000)\n", totalSetelahHapus)

	// 6. Uji Serialisasi JSON (Pastikan InternalPromoCode TIDAK MUNCUL!)
	jsonBytes, _ := json.Marshal(cart)
	fmt.Println("8. Output JSON Cart:")
	fmt.Println(string(jsonBytes))
}
