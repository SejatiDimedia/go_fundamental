package domain

import (
	"fmt"
	"time"
)

type BaseModel struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
}

type Product struct {
	BaseModel
	Name  string `json:"name"`
	Price int    `json:"price"`
	Stock int    `json:"stock"`
}

func NewProduct(id, name string, price, stock int) (*Product, error) {
	if price < 0 {
		return nil, fmt.Errorf("harga produk %s tidak boleh negatif", name)
	}
	if stock < 0 {
		return nil, fmt.Errorf("stock awal produk %s tidak boleh negatif", name)
	}

	return &Product{
		BaseModel: BaseModel{
			ID:        id,
			CreatedAt: time.Now(),
		},
		Name:  name,
		Price: price,
		Stock: stock,
	}, nil
}

func (p *Product) ReduceStock(qty int) error {
	if qty <= 0 {
		return ErrInvalidQuantity
	}

	if p.Stock < qty {
		return fmt.Errorf("gagal mengurangi stock %s (sisa: %d, diminta: %d): %w", p.Name, p.Stock, qty, ErrOutOfStock)
	}

	p.Stock -= qty
	return nil
}

func (p *Product) AddStock(qty int) error {
	if qty <= 0 {
		return ErrInvalidQuantity
	}

	p.Stock += qty
	return nil
}
