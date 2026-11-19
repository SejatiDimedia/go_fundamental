package main

import "fmt"

type User struct {
	Name string
}

func (u *User) ChangeName(name string) {
	u.Name = name
}

func main() {
	// SOLUSI: Simpan POINTER ke struct (*User) di dalam map, bukan struct value biasa!
	users := map[string]*User{
		"admin": &User{Name: "Alice"},
	}

	fmt.Println("Sebelum ChangeName:", users["admin"].Name)

	// Sekarang BISA dipanggil tanpa compile error!
	users["admin"].ChangeName("Bob")

	fmt.Println("Setelah ChangeName :", users["admin"].Name)
}
