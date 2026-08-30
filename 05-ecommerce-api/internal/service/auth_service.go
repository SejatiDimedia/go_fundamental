package service

import (
	"errors"
	"fmt"
	"time"

	"go_fundamental/05-ecommerce-api/internal/model"
	"go_fundamental/05-ecommerce-api/internal/repository"

	"github.com/golang-jwt/jwt/v5"
)

var JwtSecret = []byte("ECOMMERCE_ULTRA_SECRET_KEY_2026")

type AuthService struct {
	store *repository.Store
}

func NewAuthService(store *repository.Store) *AuthService {
	return &AuthService{store: store}
}

type JWTClaims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func (s *AuthService) Register(fullName, email, password string) (*model.User, error) {
	if _, exists := s.store.FindUserByEmail(email); exists {
		return nil, errors.New("email sudah terdaftar")
	}

	user := model.User{
		ID:       fmt.Sprintf("USR-%d", time.Now().UnixNano()%10000),
		Email:    email,
		Password: password,
		FullName: fullName,
		Role:     "customer",
	}

	s.store.SaveUser(user)
	return &user, nil
}

func (s *AuthService) Login(email, password string) (string, *model.User, error) {
	user, exists := s.store.FindUserByEmail(email)
	if !exists || user.Password != password {
		return "", nil, errors.New("email atau password salah")
	}

	claims := JWTClaims{
		UserID: user.ID,
		Email:  user.Email,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(4 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString(JwtSecret)
	if err != nil {
		return "", nil, err
	}

	return tokenStr, &user, nil
}
