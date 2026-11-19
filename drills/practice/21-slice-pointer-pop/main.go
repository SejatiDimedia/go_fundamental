package main

import "fmt"

// Pop mengambil elemen terakhir dari slice dan menghapusnya dari slice asli pemanggil.
func Pop(s *[]int) int {
	// TODO: Lengkapi fungsi ini
	lastIndex := len(*s) - 1     // 2
	lastValue := (*s)[lastIndex] // 30
	*s = (*s)[:lastIndex]        // [10 20]
	return lastValue             // 30
}

func main() {
	scores := []int{10, 20, 30}
	fmt.Println("Sebelum Pop:", scores)

	last := Pop(&scores)
	fmt.Println("Elemen yang di-Pop:", last) // Ekspektasi: 30
	fmt.Println("Setelah Pop:", scores)      // Ekspektasi: [10 20]
}
