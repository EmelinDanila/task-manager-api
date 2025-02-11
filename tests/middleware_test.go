package tests

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/EmelinDanila/task-manager-api/middleware"
	"github.com/EmelinDanila/task-manager-api/services"
	"github.com/stretchr/testify/assert"
)

func TestAuthMiddleware_ValidToken(t *testing.T) {
	authService := services.NewAuthService()
	userID := uint(123)

	// Generate a valid token
	token, _ := authService.GenerateToken(userID)

	// Create middleware and wrap the test handler
	handler := middleware.AuthMiddleware(authService)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		retrievedUserID, ok := middleware.GetUserID(r.Context())
		assert.True(t, ok, "UserID should be present in context")
		assert.Equal(t, userID, retrievedUserID, "UserID should match the one in the token")
		w.WriteHeader(http.StatusOK)
	}))

	// Create a request with a valid token
	req := httptest.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()

	// Execute the request
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code, "Valid token should allow access")
}

func TestAuthMiddleware_InvalidToken(t *testing.T) {
	authService := services.NewAuthService()

	// Create middleware and wrap the test handler
	handler := middleware.AuthMiddleware(authService)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	// Create a request with an invalid token
	req := httptest.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer invalid_token")
	rr := httptest.NewRecorder()

	// Execute the request
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code, "Invalid token should return 401")
}

func TestAuthMiddleware_MissingToken(t *testing.T) {
	authService := services.NewAuthService()

	// Create middleware and wrap the test handler
	handler := middleware.AuthMiddleware(authService)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	// Create a request without Authorization header
	req := httptest.NewRequest("GET", "/protected", nil)
	rr := httptest.NewRecorder()

	// Execute the request
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code, "Missing token should return 401")
}

func TestAuthMiddleware_InvalidHeaderFormat(t *testing.T) {
	authService := services.NewAuthService()

	// Create middleware and wrap the test handler
	handler := middleware.AuthMiddleware(authService)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	// Create a request with an invalid Authorization header format
	req := httptest.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "InvalidFormat")
	rr := httptest.NewRecorder()

	// Execute the request
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code, "Invalid header format should return 401")
}

func TestGetUserID(t *testing.T) {
	userID := uint(123)
	ctx := context.WithValue(context.Background(), middleware.UserIDKey, userID)

	retrievedUserID, ok := middleware.GetUserID(ctx)
	assert.True(t, ok, "UserID should be retrievable from context")
	assert.Equal(t, userID, retrievedUserID, "Retrieved UserID should match the original")
}

func TestGetUserID_Missing(t *testing.T) {
	ctx := context.Background() // Context without UserID

	retrievedUserID, ok := middleware.GetUserID(ctx)
	assert.False(t, ok, "UserID should not be present in context")
	assert.Equal(t, uint(0), retrievedUserID, "Default UserID should be 0 when not present")
}
