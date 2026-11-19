package main

import "fmt"

func main() {
	var input any = "100" // String!

	fmt.Println("Mencoba Kasus 1...")
	// Kasus 1: Tanpa comma-ok
	val1 := input.(int) // Ini akan PANIC!
	fmt.Println("Hasil:", val1)
}
