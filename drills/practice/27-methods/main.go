package main

import "fmt"

type User struct {
	Name string
}

// 1. Method dengan Value Receiver (u User)
func (u User) ChangeNameValue(newName string) {
	u.Name = newName
	fmt.Printf("  [Inside ChangeNameValue] Address of u: %p, Name: %s\n", &u, u.Name)
}

// 2. Method dengan Pointer Receiver (u *User)
func (u *User) ChangeNamePointer(newName string) {
	u.Name = newName
	fmt.Printf("  [Inside ChangeNamePointer] Address of u: %p, Name: %s\n", u, u.Name)
}

func main() {
	user := User{Name: "Alice"}
	fmt.Printf("Original user address: %p, Name: %s\n\n", &user, user.Name)

	fmt.Println("--- 1. Panggil user.ChangeNameValue(\"Bob\") ---")
	user.ChangeNameValue("Bob")
	fmt.Printf("Hasil di main setelah ChangeNameValue: %s\n\n", user.Name)

	fmt.Println("--- 2. Panggil user.ChangeNamePointer(\"Bob\") ---")
	user.ChangeNamePointer("Bob")
	fmt.Printf("Hasil di main setelah ChangeNamePointer: %s\n", user.Name)
}
