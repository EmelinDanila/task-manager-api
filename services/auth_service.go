package services

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// AuthService defines the interface for authentication-related operations.
type AuthService interface {
	GenerateToken(userID uint) (string, error)         // Generate a JWT for a given user ID.
	VerifyToken(tokenString string) (uint, error)      // Verify a JWT and return the user ID.
	ParseToken(tokenString string) (*jwt.Token, error) // Optionally parse token for advanced use cases.
}

type authService struct {
	secretKey string
}

// NewAuthService creates a new instance of AuthService.
func NewAuthService() AuthService {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "default_secret" // Fallback for testing (avoid in production)
	}
	return &authService{secretKey: secret}
}

// GenerateToken generates a JWT for the given user ID.
func (a *authService) GenerateToken(userID uint) (string, error) {
	claims := jwt.MapClaims{
		"userID": userID,
		"exp":    time.Now().Add(time.Hour * 24).Unix(), // 24-hour expiration
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(a.secretKey))
}

// VerifyToken validates a JWT and extracts the user ID.
func (a *authService) VerifyToken(tokenString string) (uint, error) {
	token, err := a.ParseToken(tokenString)
	if err != nil {
		return 0, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if userID, ok := claims["userID"].(float64); ok {
			return uint(userID), nil
		}
	}
	return 0, errors.New("invalid token claims")
}

// ParseToken parses and validates a JWT, returning the token for advanced use cases.
func (a *authService) ParseToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(a.secretKey), nil
	})
}
