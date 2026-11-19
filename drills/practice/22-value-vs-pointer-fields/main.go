package main

import "fmt"

// Kasus A: Field berupa VALUE (Bukan Pointer)
type CustomerValue struct {
	Name string
}

type OrderWithValue struct {
	ID       string
	Customer CustomerValue // Struct biasa, BUKAN pointer
}

// Kasus B: Field berupa POINTER
type CustomerPtr struct {
	Name string
}

type OrderWithPtr struct {
	ID       string
	Customer *CustomerPtr // Pointer
}

func main() {
	fmt.Println("=== Kasus A: Field Berupa VALUE ===")
	o1 := OrderWithValue{
		ID:       "ORD-A1",
		Customer: CustomerValue{Name: "Alice"},
	}
	o2 := o1 // Otomatis copy seluruh struct (terpisah mandiri di RAM)
	o2.Customer.Name = "Bob"

	fmt.Printf("o1.Customer.Name = %s (Alamat: %p)\n", o1.Customer.Name, &o1.Customer)
	fmt.Printf("o2.Customer.Name = %s (Alamat: %p)\n\n", o2.Customer.Name, &o2.Customer)

	fmt.Println("=== Kasus B: Field Berupa POINTER ===")
	p1 := OrderWithPtr{
		ID:       "ORD-B1",
		Customer: &CustomerPtr{Name: "Alice"},
	}
	p2 := p1 // Shallow copy: menyalin alamat pointer Customer
	// Agar terpisah, harus dialokasikan baru dengan &CustomerPtr{}
	p2.Customer = &CustomerPtr{Name: p1.Customer.Name}
	p2.Customer.Name = "Bob"

	fmt.Printf("p1.Customer.Name = %s (Alamat: %p)\n", p1.Customer.Name, p1.Customer)
	fmt.Printf("p2.Customer.Name = %s (Alamat: %p)\n", p2.Customer.Name, p2.Customer)
}
