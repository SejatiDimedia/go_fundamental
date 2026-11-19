package main

import "fmt"

func main() {
	fmt.Println("=== 1. UJI LABELED BREAK ===")
	// TODO 1:
	// Buat perulangan 'for' tak hingga (for { ... })
	// Di dalamnya ada switch status.
	// Gunakan LABELED BREAK agar ketika status == "STOP",
	// program benar-benar keluar dari perulangan for, BUKAN cuma keluar dari switch!

	commands := []string{"CONTINUE", "CONTINUE", "STOP", "CONTINUE"}
	idx := 0
MainLoop:

	// Beri nama label di sini (misal: CommandLoop:)
	for {
		cmd := commands[idx]
		idx++
		switch cmd {
		case "CONTINUE":
			fmt.Println("Memproses...")
		case "STOP":
			fmt.Println("Perintah STOP diterima! Keluar...")
			// break label di sini
			break MainLoop
		}
	}
	fmt.Println("Berhasil keluar dari loop!")

	fmt.Println("\n=== 2. UJI TAGLESS SWITCH (PENGGANTI IF-ELSE) ===")
	// TODO 2:
	// Implementasikan fungsi evaluasi nilai skor:
	// - score >= 90: "A"
	// - score >= 80: "B"
	// - score >= 70: "C"
	// - lainnya    : "D"
	// Wajib gunakan: switch tanpa variabel (switch { case ... })
	score := 85
	var grade string
	// Tulis switch di sini:
	switch {
	case score >= 90:
		grade = "A"
	case score >= 80:
		grade = "B"
	case score >= 70:
		grade = "C"
	default:
		grade = "D"
	}

	fmt.Printf("Skor %d mendapat grade: %s (Ekspektasi: B)\n", score, grade)

	fmt.Println("\n=== 3. UJI LOOP VARIABLE (Go 1.22+) ===")
	// Kita buktikan apakah kompilermu mengalokasikan alamat baru di tiap iterasi:
	values := []int{10, 20, 30}
	var ptrs []*int

	for _, v := range values {
		ptrs = append(ptrs, &v)
	}

	fmt.Print("Isi ptrs: ")
	for _, p := range ptrs {
		fmt.Printf("%d (alamat: %p) ", *p, p)
	}
	fmt.Println()

	fmt.Println("\n=== 4. UJI ITERASI STRING: BYTE VS RUNE ===")
	text := "Go👍Mantap"
	fmt.Println("Panjang len(text):", len(text), "bytes")

	fmt.Println("\n--- CARA A (Loop Byte by Byte) ---")
	for i := 0; i < len(text); i++ {
		fmt.Printf("Index byte %d: %c (hex: %x)\n", i, text[i], text[i])
	}

	fmt.Println("\n--- CARA B (Loop Rune by Rune via range) ---")
	for idx, r := range text {
		fmt.Printf("Byte offset %d: %c\n", idx, r)
	}
}
