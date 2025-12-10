package middleware

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/ellismcdougald/edmonton-bike-map/internal/utils"
)

type contextKey string

const userIDKey contextKey = "userID"

// UserIDToContext is a test helper that adds a user ID to the context.
func UserIDToContext(ctx context.Context, userID int64) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}

// AuthMiddleware validates the JWT token in the Authorization header
// and stores the user ID in the request context.
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing authorization header", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
			return
		}

		claims, err := utils.ValidateJWT(tokenString)
		if err != nil {
			log.Printf("auth_middleware: token validation failed: %v", err)
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), userIDKey, claims.UserID)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

// UserIDFromContext retrieves the authenticated user ID from the context.
// Returns the ID and true if it exists, 0 and false otherwise.
func UserIDFromContext(ctx context.Context) (int64, bool) {
	id, ok := ctx.Value(userIDKey).(int64)
	return id, ok
}
