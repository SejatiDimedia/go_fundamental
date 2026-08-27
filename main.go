package main

import (
	"fmt"
	"strings"
)

func main() {
	// ==========================================
	// 1. Membersihkan Input Pengguna (Sanitasi)
	// ==========================================
	fmt.Println("=== 1. Sanitasi String ===")
	rawUserInput := "   budi.santoso@Email.COM \n\t"
	// Menghapus whitespace di awal/akhir dan mengubah ke huruf kecil
	cleanedEmail := strings.ToLower(strings.TrimSpace(rawUserInput))
	fmt.Printf("Input mentah : %q\n", rawUserInput)
	fmt.Printf("Email bersih : %q\n", cleanedEmail)
	// ==========================================
	// 2. Pencarian & Validasi Prefix/Suffix
	// ==========================================
	fmt.Println("\n=== 2. Validasi String ===")
	invoiceNumber := "INV-2026-X8910"
	isInvoice := strings.HasPrefix(invoiceNumber, "INV-")
	isPDF := strings.HasSuffix(invoiceNumber, ".pdf")
	containsYear := strings.Contains(invoiceNumber, "2026")
	fmt.Printf("Nomor: %s\n", invoiceNumber)
	fmt.Printf("Apakah berawalan 'INV-'? %t\n", isInvoice)
	fmt.Printf("Apakah berakhiran '.pdf'? %t\n", isPDF)
	fmt.Printf("Mengandung tahun 2026? %t\n", containsYear)
	// ==========================================
	// 3. Split & Join
	// ==========================================
	fmt.Println("\n=== 3. Split & Join ===")
	tagsInput := "golang,backend,microservices,payment"
	// Memecah menjadi slice
	tagList := strings.Split(tagsInput, ",")
	fmt.Printf("Hasil Split (Slice): %#v (Total: %d tags)\n", tagList, len(tagList))
	// Menggabungkan kembali dengan pemisah baru
	hashtagText := "#" + strings.Join(tagList, " #")
	fmt.Printf("Hasil Join: %s\n", hashtagText)
	// ==========================================
	// 4. Sprintf (Membuat formatted string tanpa mencetak langsung)
	// ==========================================
	fmt.Println("\n=== 4. fmt.Sprintf ===")
	orderID := 1204
	buyer := "Andi"
	amount := 450000.75
	// Sprintf menghasilkan string baru yang bisa disimpan ke variabel / database
	message := fmt.Sprintf("Halo %s, pesanan #%d sejumlah Rp%.2f telah diterima!", buyer, orderID, amount)
	fmt.Println("Pesan:", message)
}
