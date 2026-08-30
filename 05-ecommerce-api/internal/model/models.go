package model

import "time"

type User struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Password string `json:"-"`
	FullName string `json:"full_name"`
	Role     string `json:"role"`
}

type Product struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       int    `json:"price"`
	Stock       int    `json:"stock"`
}

type CartItem struct {
	Product  Product `json:"product"`
	Quantity int     `json:"quantity"`
}

type Order struct {
	ID          string     `json:"id"`
	UserID      string     `json:"user_id"`
	Items       []CartItem `json:"items"`
	TotalAmount int        `json:"total_amount"`
	Status      string     `json:"status"` // "PAID", "FAILED"
	CreatedAt   time.Time  `json:"created_at"`
}

type StandardResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
	Error   any    `json:"error,omitempty"`
}
