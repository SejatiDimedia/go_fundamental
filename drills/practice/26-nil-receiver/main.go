package main

import "fmt"

type Node struct {
	Value int
	Next  *Node
}

// Method dengan pointer receiver
func (n *Node) GetValue() int {
	// Menghandle nil receiver dengan aman
	if n == nil {
		return 0
	}
	return n.Value
}

func main() {
	var head *Node // nil pointer
	fmt.Printf("head ptr: %v\n", head)

	// Memanggil method pada nil pointer:
	result := head.GetValue()
	fmt.Printf("Hasil head.GetValue(): %d\n", result)

	// Jika ada nilainya:
	activeNode := &Node{Value: 42}
	fmt.Printf("Hasil activeNode.GetValue(): %d\n", activeNode.GetValue())
}
