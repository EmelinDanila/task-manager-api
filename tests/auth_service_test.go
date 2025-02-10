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
	os.Setenv("JWT_SECRET", "test_secret") // Set a temporary secret key for testing.
	authService := services.NewAuthService()

	// Test token generation.
	userID := uint(123)
	token, err := authService.GenerateToken(userID)
	assert.NoError(t, err)
	assert.NotEmpty(t, token, "Token should not be empty")
}

func TestAuthService_VerifyToken(t *testing.T) {
	os.Setenv("JWT_SECRET", "test_secret") // Set a temporary secret key for testing.
	authService := services.NewAuthService()

	// Generate and verify the token.
	userID := uint(123)
	token, err := authService.GenerateToken(userID)
	assert.NoError(t, err)

	// Validate the token.
	parsedUserID, err := authService.VerifyToken(token)
	assert.NoError(t, err)
	assert.Equal(t, userID, parsedUserID, "Parsed userID should match the original")
}

func TestAuthService_ParseToken(t *testing.T) {
	os.Setenv("JWT_SECRET", "test_secret") // Set a temporary secret key for testing.
	authService := services.NewAuthService()

	// Generate the token.
	userID := uint(123)
	token, err := authService.GenerateToken(userID)
	assert.NoError(t, err)

	// Parse the token.
	parsedToken, err := authService.ParseToken(token)
	assert.NoError(t, err)
	assert.NotNil(t, parsedToken, "Parsed token should not be nil")

	// Verify claims.
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

	// Attempt to parse an invalid token.
	_, err := authService.VerifyToken("invalid_token")
	assert.Error(t, err, "Invalid token should return an error")
}
