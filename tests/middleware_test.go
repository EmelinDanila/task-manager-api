package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/EmelinDanila/task-manager-api/middleware"
	"github.com/EmelinDanila/task-manager-api/services"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestAuthMiddleware_ValidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	authService := services.NewAuthService()
	userID := uint(123)

	// Генерируем валидный токен
	token, _ := authService.GenerateToken(userID)

	// Создаем тестовый HTTP запрос
	req := httptest.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	// Создаем ResponseRecorder для теста
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Вызываем middleware
	middleware.AuthMiddleware(authService)(c)

	// Проверяем статус-код
	assert.Equal(t, http.StatusOK, w.Code, "Valid token should allow access")
}

func TestAuthMiddleware_InvalidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	authService := services.NewAuthService()

	req := httptest.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer invalid_token")
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	middleware.AuthMiddleware(authService)(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code, "Invalid token should return 401")
}

func TestAuthMiddleware_MissingToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	authService := services.NewAuthService()

	req := httptest.NewRequest("GET", "/protected", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	middleware.AuthMiddleware(authService)(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code, "Missing token should return 401")
}

func TestAuthMiddleware_InvalidHeaderFormat(t *testing.T) {
	gin.SetMode(gin.TestMode)
	authService := services.NewAuthService()

	req := httptest.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "InvalidFormat")
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	middleware.AuthMiddleware(authService)(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code, "Invalid header format should return 401")
}

func TestGetUserID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	userID := uint(123)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Set("userID", userID)

	retrievedUserID, ok := middleware.GetUserID(c)
	assert.True(t, ok, "UserID should be retrievable from context")
	assert.Equal(t, userID, retrievedUserID, "Retrieved UserID should match the original")
}

func TestGetUserID_Missing(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder()) // Context без UserID

	retrievedUserID, ok := middleware.GetUserID(c)
	assert.False(t, ok, "UserID should not be present in context")
	assert.Equal(t, uint(0), retrievedUserID, "Default UserID should be 0 when not present")
}
