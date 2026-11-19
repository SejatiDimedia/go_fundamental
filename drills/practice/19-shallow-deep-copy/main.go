package main

import "fmt"

type Customer struct {
	Name string
}

type Order struct {
	ID       string
	Customer *Customer
}

func main() {
	cust := &Customer{Name: "Alice"}

	order1 := Order{
		ID:       "ORD-001",
		Customer: cust,
	}

	fmt.Printf("Alamat Customer asli: %p\n", cust)
	fmt.Printf("order1: ID=%s, CustomerPtr=%p, Name=%s\n\n", order1.ID, order1.Customer, order1.Customer.Name)

	// Shallow copy
	order2 := order1
	fmt.Printf("order2 (setelah copy): ID=%s, CustomerPtr=%p, Name=%s\n\n", order2.ID, order2.Customer, order2.Customer.Name)

	// Modifikasi order2
	order2.ID = "ORD-002"
	order2.Customer.Name = "Bob"

	fmt.Println("--- Hasil Setelah order2 Dimodifikasi ---")
	fmt.Println("1. order1.ID =", order1.ID)
	fmt.Println("2. order1.Customer.Name =", order1.Customer.Name)
	fmt.Println("   order2.ID =", order2.ID)
	fmt.Println("   order2.Customer.Name =", order2.Customer.Name)
}
