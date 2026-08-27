package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"go_fundamental/02-inventory-management/domain"
	"go_fundamental/02-inventory-management/repository"
)

func readLine(scanner *bufio.Scanner, label string) string {
	fmt.Print(label)
	scanner.Scan()
	return strings.TrimSpace(scanner.Text())
}

func seedInitialData(store *repository.MemoryStore[*domain.Product]) {
	p1, _ := domain.NewProduct("PRD-001", "Laptop ThinkPad", 15000000, 5)
	p2, _ := domain.NewProduct("PRD-002", "Logitech MX Master Mouse", 1200000, 10)
	p3, _ := domain.NewProduct("PRD-003", "Keychron Keyboard", 1500000, 8)
	p4, _ := domain.NewProduct("PRD-004", "USB-C Hub Multiport", 350000, 2)

	store.Save(p1.ID, p1)
	store.Save(p2.ID, p2)
	store.Save(p3.ID, p3)
	store.Save(p4.ID, p4)
}

func displayInventory(store *repository.MemoryStore[*domain.Product]) {
	products := store.GetAll()

	fmt.Println("\n========================= DAFTAR INVENTARIS TOKO =========================")
	fmt.Printf("%-10s | %-26s | %-14s | %-8s | %-16s\n", "ID Produk", "Nama Barang", "Harga", "Stok", "Status")
	fmt.Println("--------------------------------------------------------------------------")

	totalAssetValue := 0
	for _, p := range products {
		status := "🟢 Tersedia"
		if p.Stock == 0 {
			status = "🔴 Habis"
		} else if p.Stock <= 3 {
			status = "🟡 Menipis"
		}

		totalAssetValue += p.Price * p.Stock
		fmt.Printf("%-10s | %-26s | Rp%-12d | %-8d | %s\n", p.ID, p.Name, p.Price, p.Stock, status)
	}
	fmt.Println("--------------------------------------------------------------------------")
	fmt.Printf("TOTAL NILAI ASET INVENTARIS: Rp%d\n", totalAssetValue)
	fmt.Println("==========================================================================")
}

func handleCheckout(scanner *bufio.Scanner, store *repository.MemoryStore[*domain.Product]) {
	displayInventory(store)

	prodID := strings.ToUpper(readLine(scanner, "\nMasukkan ID Produk yang dibeli: "))
	product, exists := store.FindByID(prodID)
	if !exists {
		fmt.Println("❌ Error: Produk tidak ditemukan!")
		return
	}

	qtyStr := readLine(scanner, fmt.Sprintf("Masukkan jumlah unit %s yang dibeli: ", product.Name))
	qty, err := strconv.Atoi(qtyStr)
	if err != nil || qty <= 0 {
		fmt.Println("❌ Error: Jumlah beli harus berupa angka positif!")
		return
	}

	// 1. Validasi pengurangan stok
	if err := product.ReduceStock(qty); err != nil {
		// Demonstrasi errors.Is untuk mendeteksi ErrOutOfStock
		if errors.Is(err, domain.ErrOutOfStock) {
			fmt.Printf("❌ Transaksi Gagal: %s\n", err.Error())
		} else {
			fmt.Println("❌ Gagal:", err)
		}
		return
	}

	totalBill := product.Price * qty
	fmt.Printf("\nTotal Tagihan: Rp%d (%d unit x %s)\n", totalBill, qty, product.Name)

	// 2. Pilih Metode Pembayaran (Interface Polymorphism)
	fmt.Println("Pilih Metode Pembayaran:")
	fmt.Println("1. Tunai (Cash)")
	fmt.Println("2. Virtual Account BCA")
	fmt.Println("3. E-Wallet (GoPay)")

	methodChoice := readLine(scanner, "Pilihan (1-3): ")
	var payment domain.PaymentMethod

	switch methodChoice {
	case "1":
		cashStr := readLine(scanner, "Masukkan jumlah uang tunai: Rp")
		cash, _ := strconv.Atoi(cashStr)
		payment = &domain.CashPayment{CashGiven: cash}

	case "2":
		payment = &domain.VAPayment{
			BankName: "BCA",
			VANumber: "880011223344",
		}

	case "3":
		payment = &domain.EWalletPayment{
			Provider: "GoPay",
			WalletID: "081299887766",
			Balance:  20000000, // Saldo simulasi Rp20 Juta
		}

	default:
		fmt.Println("❌ Pilihan tidak valid, pesanan dibatalkan. Mengembalikan stok...")
		product.AddStock(qty) // Rollback stok
		return
	}

	// 3. Proses Pembayaran via Interface
	ref, payErr := payment.Pay(totalBill)
	if payErr != nil {
		fmt.Printf("❌ Pembayaran Gagal: %s\n", payErr.Error())
		fmt.Println(">> Mengembalikan stok ke inventaris (Rollback)...")
		product.AddStock(qty)
		return
	}

	fmt.Println("\n================ STRUK PEMBELIAN ================")
	fmt.Printf("Metode   : %s\n", payment.Name())
	fmt.Printf("Produk   : %s (%d unit)\n", product.Name, qty)
	fmt.Printf("Total    : Rp%d\n", totalBill)
	fmt.Printf("Status   : LUNAS (%s)\n", ref)
	fmt.Println("=================================================")
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	productStore := repository.NewMemoryStore[*domain.Product]()

	// Mengisi data awal
	seedInitialData(productStore)

	fmt.Println("=========================================================")
	fmt.Println("   🏬 SISTEM MANAJEMEN INVENTARIS & KASIR MODERN (Go)   ")
	fmt.Println("=========================================================")

	for {
		fmt.Println("\nMenu Utama:")
		fmt.Println("1. Lihat Daftar Inventaris & Status Stok")
		fmt.Println("2. Tambah Produk Baru")
		fmt.Println("3. Tambah Stok Barang (Restock)")
		fmt.Println("4. Transaksi Penjualan (Checkout)")
		fmt.Println("5. Keluar")

		choice := readLine(scanner, "Pilih menu (1-5): ")

		switch choice {
		case "1":
			displayInventory(productStore)

		case "2":
			id := strings.ToUpper(readLine(scanner, "ID Produk baru (contoh: PRD-005): "))
			name := readLine(scanner, "Nama Produk: ")
			priceStr := readLine(scanner, "Harga Produk (Rp): ")
			stockStr := readLine(scanner, "Stok Awal: ")

			price, _ := strconv.Atoi(priceStr)
			stock, _ := strconv.Atoi(stockStr)

			newProd, err := domain.NewProduct(id, name, price, stock)
			if err != nil {
				fmt.Printf("❌ Gagal membuat produk: %s\n", err.Error())
				continue
			}

			productStore.Save(newProd.ID, newProd)
			fmt.Printf("✅ Produk %s berhasil ditambahkan!\n", newProd.Name)

		case "3":
			displayInventory(productStore)
			id := strings.ToUpper(readLine(scanner, "Masukkan ID Produk yang akan di-restock: "))
			product, exists := productStore.FindByID(id)
			if !exists {
				fmt.Println("❌ Produk tidak ditemukan!")
				continue
			}

			qtyStr := readLine(scanner, fmt.Sprintf("Tambah berapa stok untuk %s? ", product.Name))
			qty, _ := strconv.Atoi(qtyStr)

			if err := product.AddStock(qty); err != nil {
				fmt.Printf("❌ Error: %s\n", err.Error())
			} else {
				fmt.Printf("✅ Stok %s sekarang menjadi %d unit.\n", product.Name, product.Stock)
			}

		case "4":
			handleCheckout(scanner, productStore)

		case "5":
			fmt.Println("Terima kasih telah menggunakan sistem inventaris. Sampai jumpa!")
			return

		default:
			fmt.Println("❌ Pilihan menu tidak valid!")
		}
	}
}
