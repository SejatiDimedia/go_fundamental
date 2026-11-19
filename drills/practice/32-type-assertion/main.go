package main

import "fmt"

func main() {
	var data any = 100

	// ❌ data + 50 tidak bisa:
	// result := data + 50 // COMPILE ERROR: operator + not defined on any

	// ✅ Cara 1: Type Assertion langsung (Hati-hati, bisa PANIC jika salah tipe!)
	val := data.(int)
	fmt.Println("Hasil Type Assertion langsung:", val+50)

	// 🛡️ Cara 2: Type Assertion Aman (Comma-Ok Idiom)
	str, ok := data.(string)
	if !ok {
		fmt.Println("Bukan string! Selamat dari panic!")
	} else {
		fmt.Println("String:", str)
	}

	// 🔀 Cara 3: Type Switch (Memeriksa banyak tipe sekaligus)
	inspectType(data)
	inspectType("Halo Dunia")
	inspectType(true)
}

func inspectType(val any) {
	switch v := val.(type) {
	case int:
		fmt.Printf("[Type Switch] Integer: %d (bisa dikali 2 = %d)\n", v, v*2)
	case string:
		fmt.Printf("[Type Switch] String: %s (panjang: %d)\n", v, len(v))
	case bool:
		fmt.Printf("[Type Switch] Boolean: %t\n", v)
	default:
		fmt.Printf("[Type Switch] Tipe lain: %T\n", v)
	}
}
