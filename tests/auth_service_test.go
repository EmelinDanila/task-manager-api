package tests

import (
	"os"
	"testing"
	"time"

	"github.com/EmelinDanila/task-manager-api/services"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestAuthService_GenerateToken(t *testing.T) {
	os.Setenv("JWT_SECRET", "test_secret") // Устанавливаем временный секретный ключ для тестов.
	authService := services.NewAuthService()

	// Тестируем создание токена.
	userID := uint(123)
	token, err := authService.GenerateToken(userID)
	assert.NoError(t, err)
	assert.NotEmpty(t, token, "Token should not be empty")
}

func TestAuthService_VerifyToken(t *testing.T) {
	os.Setenv("JWT_SECRET", "test_secret") // Устанавливаем временный секретный ключ для тестов.
	authService := services.NewAuthService()

	// Создаем токен и проверяем его.
	userID := uint(123)
	token, err := authService.GenerateToken(userID)
	assert.NoError(t, err)

	// Проверяем валидность токена.
	parsedUserID, err := authService.VerifyToken(token)
	assert.NoError(t, err)
	assert.Equal(t, userID, parsedUserID, "Parsed userID should match the original")
}

func TestAuthService_ParseToken(t *testing.T) {
	os.Setenv("JWT_SECRET", "test_secret") // Устанавливаем временный секретный ключ для тестов.
	authService := services.NewAuthService()

	// Создаем токен.
	userID := uint(123)
	token, err := authService.GenerateToken(userID)
	assert.NoError(t, err)

	// Парсим токен.
	parsedToken, err := authService.ParseToken(token)
	assert.NoError(t, err)
	assert.NotNil(t, parsedToken, "Parsed token should not be nil")

	// Проверяем claims.
	if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok && parsedToken.Valid {
		assert.Equal(t, float64(userID), claims["userID"], "Claims should contain the correct userID")
		assert.True(t, claims["exp"].(float64) > float64(time.Now().Unix()), "Token should not be expired")
	} else {
		t.Errorf("Failed to parse claims")
	}
}

func TestAuthService_InvalidToken(t *testing.T) {
	os.Setenv("JWT_SECRET", "test_secret")
	authService := services.NewAuthService()

	// Пробуем парсить недействительный токен.
	_, err := authService.VerifyToken("invalid_token")
	assert.Error(t, err, "Invalid token should return an error")
}
