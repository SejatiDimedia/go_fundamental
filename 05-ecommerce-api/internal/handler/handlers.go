package handler

import (
	"net/http"

	"go_fundamental/05-ecommerce-api/internal/model"
	"go_fundamental/05-ecommerce-api/internal/repository"
	"go_fundamental/05-ecommerce-api/internal/service"

	"github.com/gin-gonic/gin"
)

type AppHandler struct {
	authService  *service.AuthService
	orderService *service.OrderService
	store        *repository.Store
}

func NewAppHandler(auth *service.AuthService, order *service.OrderService, store *repository.Store) *AppHandler {
	return &AppHandler{
		authService:  auth,
		orderService: order,
		store:        store,
	}
}

// Auth Handlers
type RegisterReq struct {
	FullName string `json:"full_name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

func (h *AppHandler) Register(c *gin.Context) {
	var req RegisterReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.StandardResponse{Success: false, Message: "Input tidak valid", Error: err.Error()})
		return
	}

	user, err := h.authService.Register(req.FullName, req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.StandardResponse{Success: false, Message: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, model.StandardResponse{Success: true, Message: "Registrasi berhasil!", Data: user})
}

type LoginReq struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (h *AppHandler) Login(c *gin.Context) {
	var req LoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.StandardResponse{Success: false, Message: "Input tidak valid", Error: err.Error()})
		return
	}

	token, user, err := h.authService.Login(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, model.StandardResponse{Success: false, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, model.StandardResponse{
		Success: true,
		Message: "Login sukses!",
		Data: gin.H{
			"token": token,
			"user":  user,
		},
	})
}

// Product Handlers
func (h *AppHandler) GetProducts(c *gin.Context) {
	products := h.store.GetAllProducts()
	c.JSON(http.StatusOK, model.StandardResponse{Success: true, Message: "Katalog produk", Data: products})
}

// Cart & Order Handlers
type AddCartReq struct {
	ProductID string `json:"product_id" binding:"required"`
	Quantity  int    `json:"quantity" binding:"required,gt=0"`
}

func (h *AppHandler) AddToCart(c *gin.Context) {
	userID := c.GetString("user_id")
	var req AddCartReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.StandardResponse{Success: false, Message: "Input tidak valid", Error: err.Error()})
		return
	}

	if err := h.orderService.AddToCart(userID, req.ProductID, req.Quantity); err != nil {
		c.JSON(http.StatusBadRequest, model.StandardResponse{Success: false, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, model.StandardResponse{Success: true, Message: "Produk berhasil ditambahkan ke keranjang"})
}

func (h *AppHandler) GetCart(c *gin.Context) {
	userID := c.GetString("user_id")
	items, total := h.orderService.GetCart(userID)

	c.JSON(http.StatusOK, model.StandardResponse{
		Success: true,
		Message: "Isi keranjang belanja",
		Data: gin.H{
			"items":        items,
			"total_amount": total,
		},
	})
}

func (h *AppHandler) Checkout(c *gin.Context) {
	userID := c.GetString("user_id")
	order, err := h.orderService.Checkout(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.StandardResponse{Success: false, Message: "Checkout gagal", Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, model.StandardResponse{Success: true, Message: "Pesanan berhasil dibuat & lunas!", Data: order})
}
