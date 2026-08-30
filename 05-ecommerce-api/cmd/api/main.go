package main

import (
	"log"
	"net/http"

	"go_fundamental/05-ecommerce-api/internal/handler"
	"go_fundamental/05-ecommerce-api/internal/middleware"
	"go_fundamental/05-ecommerce-api/internal/model"
	"go_fundamental/05-ecommerce-api/internal/repository"
	"go_fundamental/05-ecommerce-api/internal/service"

	"github.com/gin-gonic/gin"
)

func main() {
	// 1. Dependency Injection (Wiring Layers Manual Standar Go)
	store := repository.NewStore()
	authService := service.NewAuthService(store)
	orderService := service.NewOrderService(store)
	appHandler := handler.NewAppHandler(authService, orderService, store)

	// 2. Setup Router
	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, model.StandardResponse{Success: true, Message: "E-Commerce API is Healthy"})
	})

	v1 := router.Group("/api/v1")
	{
		// Public Auth
		v1.POST("/auth/register", appHandler.Register)
		v1.POST("/auth/login", appHandler.Login)

		// Public Products
		v1.GET("/products", appHandler.GetProducts)

		// Protected Cart & Checkout (Wajib Login / JWT)
		protected := v1.Group("")
		protected.Use(middleware.AuthMiddleware())
		{
			protected.POST("/cart/items", appHandler.AddToCart)
			protected.GET("/cart", appHandler.GetCart)
			protected.POST("/checkout", appHandler.Checkout)
		}
	}

	log.Println("🚀 Server E-Commerce berjalan di http://localhost:8080")
	router.Run(":8080")
}
