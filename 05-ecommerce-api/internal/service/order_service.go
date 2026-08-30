package service

import (
	"errors"
	"fmt"
	"time"

	"go_fundamental/05-ecommerce-api/internal/model"
	"go_fundamental/05-ecommerce-api/internal/repository"
)

type OrderService struct {
	store *repository.Store
}

func NewOrderService(store *repository.Store) *OrderService {
	return &OrderService{store: store}
}

func (s *OrderService) AddToCart(userID, productID string, qty int) error {
	if qty <= 0 {
		return errors.New("qty harus minimal 1")
	}

	product, exists := s.store.FindProductByID(productID)
	if !exists {
		return repository.ErrNotFound
	}

	if product.Stock < qty {
		return repository.ErrOutOfStock
	}

	item := model.CartItem{
		Product:  *product,
		Quantity: qty,
	}

	s.store.AddToCart(userID, item)
	return nil
}

func (s *OrderService) GetCart(userID string) ([]model.CartItem, int) {
	items := s.store.GetCart(userID)
	total := 0
	for _, it := range items {
		total += it.Product.Price * it.Quantity
	}
	return items, total
}

func (s *OrderService) Checkout(userID string) (*model.Order, error) {
	items, total := s.GetCart(userID)
	if len(items) == 0 {
		return nil, errors.New("keranjang belanja kosong")
	}

	order := model.Order{
		ID:          fmt.Sprintf("ORD-%d", time.Now().Unix()%1000000),
		UserID:      userID,
		Items:       items,
		TotalAmount: total,
		Status:      "PAID",
		CreatedAt:   time.Now(),
	}

	if err := s.store.CreateOrder(order); err != nil {
		return nil, err
	}

	return &order, nil
}
