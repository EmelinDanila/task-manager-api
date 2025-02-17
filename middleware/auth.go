package middleware

import (
	"net/http"
	"strings"

	"github.com/EmelinDanila/task-manager-api/services"
	"github.com/gin-gonic/gin"
)

// AuthMiddleware creates a Gin middleware for authentication
func AuthMiddleware(authService services.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract the Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
			c.Abort()
			return
		}

		// Remove "Bearer " prefix from the token
		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == authHeader {
			// If the token is the same as the header, "Bearer " prefix was not present
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization header format"})
			c.Abort()
			return
		}

		// Verify the token
		userID, err := authService.VerifyToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// Set the user ID in the context for later use
		c.Set("userID", userID)
		c.Next()
	}
}

// GetUserID retrieves the user ID from the Gin context
func GetUserID(c *gin.Context) (uint, bool) {
	userID, exists := c.Get("userID")
	if !exists {
		return 0, false
	}
	return userID.(uint), true
}
