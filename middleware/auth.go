package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/EmelinDanila/task-manager-api/services"
)

type key string

const UserIDKey key = "userID"

// AuthMiddleware checks for the presence and validity of the JWT token.
func AuthMiddleware(authService services.AuthService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Extracting the token from the Authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Authorization header missing", http.StatusUnauthorized)
				return
			}

			// Checking the header format
			token := strings.TrimPrefix(authHeader, "Bearer ")
			if token == authHeader { // If "Bearer" is not found
				http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
				return
			}

			// Verifying the token
			userID, err := authService.VerifyToken(token)
			if err != nil {
				http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
				return
			}

			// Adding userID to the context
			ctx := context.WithValue(r.Context(), UserIDKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetUserID extracts userID from the context.
func GetUserID(ctx context.Context) (uint, bool) {
	userID, ok := ctx.Value(UserIDKey).(uint)
	return userID, ok
}
