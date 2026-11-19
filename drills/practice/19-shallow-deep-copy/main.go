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
	fmt.Println("=== 1. SKENARIO SHALLOW COPY (BUG) ===")
	cust1 := &Customer{Name: "Alice"}
	order1 := Order{ID: "ORD-001", Customer: cust1}
	order2 := order1 // Shallow copy: menyalin alamat pointer cust1

	order2.ID = "ORD-002"
	order2.Customer.Name = "Bob" // Merusak order1!

	fmt.Printf("order1 (Customer Addr: %p): ID=%s, Customer.Name=%s\n", order1.Customer, order1.ID, order1.Customer.Name)
	fmt.Printf("order2 (Customer Addr: %p): ID=%s, Customer.Name=%s\n\n", order2.Customer, order2.ID, order2.Customer.Name)

	fmt.Println("=== 2. SKENARIO DEEP COPY (SOLUSI) ===")
	custOriginal := &Customer{Name: "Alice"}
	orderOrig := Order{ID: "ORD-001", Customer: custOriginal}

	// Deep copy: alokasikan struct Customer baru dengan &Customer{}
	orderCloned := orderOrig
	orderCloned.Customer = &Customer{
		Name: orderOrig.Customer.Name, // Duplikasi nilai string-nya ke rumah baru
	}

	orderCloned.ID = "ORD-002"
	orderCloned.Customer.Name = "Bob" // Aman, tidak merusak orderOrig!

	fmt.Printf("orderOrig   (Customer Addr: %p): ID=%s, Customer.Name=%s\n", orderOrig.Customer, orderOrig.ID, orderOrig.Customer.Name)
	fmt.Printf("orderCloned (Customer Addr: %p): ID=%s, Customer.Name=%s\n", orderCloned.Customer, orderCloned.ID, orderCloned.Customer.Name)
}
