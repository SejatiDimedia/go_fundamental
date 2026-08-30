package repository

import (
	"errors"
	"sync"

	"go_fundamental/05-ecommerce-api/internal/model"
)

var (
	ErrNotFound   = errors.New("data tidak ditemukan")
	ErrOutOfStock = errors.New("stok produk tidak mencukupi")
)

type Store struct {
	mu       sync.RWMutex
	Users    map[string]model.User       // Email -> User
	Products map[string]*model.Product   // ProductID -> Product
	Carts    map[string][]model.CartItem // UserID -> CartItems
	Orders   map[string]model.Order      // OrderID -> Order
}

func NewStore() *Store {
	s := &Store{
		Users:    make(map[string]model.User),
		Products: make(map[string]*model.Product),
		Carts:    make(map[string][]model.CartItem),
		Orders:   make(map[string]model.Order),
	}

	// Seed Data Produk Awal
	s.Products["PRD-101"] = &model.Product{ID: "PRD-101", Name: "Keychron Mechanical Keyboard", Price: 1250000, Stock: 10, Description: "Wireless Mechanical Keyboard RGB"}
	s.Products["PRD-102"] = &model.Product{ID: "PRD-102", Name: "Logitech MX Master 3S", Price: 1450000, Stock: 5, Description: "Ergonomic Performance Wireless Mouse"}
	s.Products["PRD-103"] = &model.Product{ID: "PRD-103", Name: "Monitor Gaming 27 Inch 165Hz", Price: 3200000, Stock: 3, Description: "IPS 2K 165Hz HDR Gaming Monitor"}

	// Seed Akun User Awal
	s.Users["budi@example.com"] = model.User{ID: "USR-001", Email: "budi@example.com", Password: "password123", FullName: "Budi Santoso", Role: "customer"}

	return s
}

// User Operations
func (s *Store) SaveUser(user model.User) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Users[user.Email] = user
}

func (s *Store) FindUserByEmail(email string) (model.User, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	u, ok := s.Users[email]
	return u, ok
}

// Product Operations
func (s *Store) GetAllProducts() []model.Product {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]model.Product, 0, len(s.Products))
	for _, p := range s.Products {
		list = append(list, *p)
	}
	return list
}

func (s *Store) FindProductByID(id string) (*model.Product, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	p, ok := s.Products[id]
	return p, ok
}

// Cart Operations
func (s *Store) GetCart(userID string) []model.CartItem {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.Carts[userID]
}

func (s *Store) AddToCart(userID string, item model.CartItem) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Carts[userID] = append(s.Carts[userID], item)
}

func (s *Store) ClearCart(userID string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.Carts, userID)
}

// Order & Stock Reduction (Atomic Transaction Simulation)
func (s *Store) CreateOrder(order model.Order) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 1. Cek ketersediaan semua stok terlebih dahulu
	for _, item := range order.Items {
		prod, exists := s.Products[item.Product.ID]
		if !exists || prod.Stock < item.Quantity {
			return ErrOutOfStock
		}
	}

	// 2. Potong stok
	for _, item := range order.Items {
		s.Products[item.Product.ID].Stock -= item.Quantity
	}

	// 3. Simpan Order & kosongkan keranjang
	s.Orders[order.ID] = order
	delete(s.Carts, order.UserID)

	return nil
}
