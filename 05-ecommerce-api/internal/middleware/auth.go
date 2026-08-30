package middleware

import (
	"net/http"
	"strings"

	"go_fundamental/05-ecommerce-api/internal/model"
	"go_fundamental/05-ecommerce-api/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, model.StandardResponse{
				Success: false,
				Message: "Token autentikasi tidak ditemukan!",
			})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.JSON(http.StatusUnauthorized, model.StandardResponse{
				Success: false,
				Message: "Format token salah! Gunakan: 'Bearer <token>'",
			})
			c.Abort()
			return
		}

		tokenString := parts[1]
		claims := &service.JWTClaims{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (any, error) {
			return service.JwtSecret, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, model.StandardResponse{
				Success: false,
				Message: "Token tidak valid atau kadaluarsa!",
			})
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)
		c.Next()
	}
}
