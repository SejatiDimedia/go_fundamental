package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type CartItem struct {
	Name     string
	Price    int
	Quantity int
}

var catalog = map[string]int{
	"kopi":   18000,
	"susu":   15000,
	"roti":   12000,
	"mie":    3500,
	"telur":  25000,
	"minyak": 32000,
}

var inventory = map[string]int{
	"kopi":   10,
	"susu":   8,
	"roti":   5,
	"mie":    20,
	"telur":  4,
	"minyak": 6,
}

func readInput(scanner *bufio.Scanner, prompt string) string {
	fmt.Print(prompt)
	scanner.Scan()
	return strings.TrimSpace(scanner.Text())
}

func displayCatalog() {
	fmt.Println("\n================ KATALOG TOKO ================")
	fmt.Printf("%-12s | %-12s | %-8s\n", "Nama Barang", "Harga (Rp)", "Sisa Stok")
	fmt.Println("----------------------------------------------")
	for item, price := range catalog {
		stock := inventory[item]
		fmt.Printf("%-12s | Rp%-10d | %-8d\n", strings.Title(item), price, stock)
	}
	fmt.Println("==============================================")
}

func displayCart(cart []CartItem) {
	fmt.Println("\n================ ISI KERANJANG ===============")
	if len(cart) == 0 {
		fmt.Println("Keranjang belanja Anda masih kosong.")
		fmt.Println("==============================================")
		return
	}
	fmt.Printf("%-4s | %-12s | %-6s | %-12s | %-12s\n", "No", "Barang", "Qty", "Harga", "Subtotal")
	fmt.Println("------------------------------------------------------------")
	total := 0
	for i, item := range cart {
		subtotal := item.Price * item.Quantity
		total += subtotal
		fmt.Printf("%-4d | %-12s | %-6d | Rp%-10d | Rp%-10d\n", i+1, strings.Title(item.Name), item.Quantity, item.Price, subtotal)
		fmt.Println("------------------------------------------------------------")
		fmt.Printf("TOTAL BELANJA: Rp%d\n", total)
		fmt.Println("============================================================")
	}

}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	var cart []CartItem

	for {
		fmt.Println("\n=== SISTEM KASIR & DAFTAR BELANJA ===")
		fmt.Println("1. Lihat Katalog Produk")
		fmt.Println("2. Tambah Barang ke Keranjang")
		fmt.Println("3. Lihat Keranjang Belanja")
		fmt.Println("4. Hapus Barang dari Keranjang")
		fmt.Println("5. Checkout & Pembayaran")
		fmt.Println("6. Keluar")

		choice := readInput(scanner, "Pilih menu (1-6): ")

		switch choice {
		case "1":
			displayCatalog()
		case "2":
			displayCatalog()
			itemName := strings.ToLower(readInput(scanner, "Masukkan nama barang yang ingin dibeli: "))

			price, exists := catalog[itemName]

			if !exists {
				fmt.Println("❌ Barang tidak ditemukan di katalog toko")
				continue
			}

			qtyStr := readInput(scanner, "Masukkan jumlah(qty): ")
			qty, err := strconv.Atoi(qtyStr)

			if err != nil || qty <= 0 {
				fmt.Println("❌ Jumlah harus berupa angka bulat positif")
				continue
			}

			if qty > inventory[itemName] {
				fmt.Printf("❌ Stok tidak mencukupi! Sisa stok %s hanya: %d\n", itemName, inventory[itemName])
				continue
			}

			inventory[itemName] -= qty
			cart = append(cart, CartItem{
				Name:     itemName,
				Price:    price,
				Quantity: qty,
			})
			fmt.Printf("✅ Berhasil menambahkan %d %s ke keranjang\n", qty, strings.Title(itemName))

		case "3":
			displayCart(cart)
		case "4":
			if len(cart) == 0 {
				fmt.Println("❌ Keranjang masih kosong!")
				continue
			}
			displayCart(cart)
			indexStr := readInput(scanner, "Masukkan nomor urut barang yang ingin dihapus: ")
			idx, err := strconv.Atoi(indexStr)
			if err != nil || idx < 1 || idx > len(cart) {
				fmt.Println("❌ Nomor urut tidak valid!")
				continue
			}
			// Mengembalikan stok toko
			removedItem := cart[idx-1]
			inventory[removedItem.Name] += removedItem.Quantity
			// Teknik menghapus elemen dari Slice di Go:
			// Menggabungkan slice sebelum index dengan slice setelah index
			cart = append(cart[:idx-1], cart[idx:]...)
			fmt.Printf("✅ Barang %s berhasil dihapus dari keranjang.\n", removedItem.Name)
		case "5":
			if len(cart) == 0 {
				fmt.Println("❌ Keranjang kosong, tidak bisa checkout!")
				continue
			}
			displayCart(cart)
			grandTotal := 0
			for _, item := range cart {
				grandTotal += item.Price * item.Quantity
			}
			fmt.Printf("\nTotal yang harus dibayar: Rp%d\n", grandTotal)
			payStr := readInput(scanner, "Masukkan nominal uang bayar: Rp")
			payAmount, err := strconv.Atoi(payStr)
			if err != nil || payAmount < grandTotal {
				fmt.Println("❌ Uang yang dibayarkan kurang atau nominal tidak valid!")
				continue
			}
			kembalian := payAmount - grandTotal
			fmt.Println("\n🎉 PEMBAYARAN BERHASIL!")
			fmt.Printf("Uang Diterima : Rp%d\n", payAmount)
			fmt.Printf("Total Belanja : Rp%d\n", grandTotal)
			fmt.Printf("Kembalian     : Rp%d\n", kembalian)
			fmt.Println("Terima kasih telah berbelanja!")
			// Reset keranjang belanja setelah berhasil checkout
			cart = make([]CartItem, 0)
		case "6":
			fmt.Println("Keluar dari program. Terima kasih!")
			return
		default:
			fmt.Println("❌ Pilihan menu tidak valid. Silakan pilih 1-6.")
		}
	}
}
