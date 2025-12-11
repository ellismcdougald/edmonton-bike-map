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

// UserIDToContext returns a copy of ctx that contains userID stored under the package's user ID context key.
// It is intended for use in tests to inject an authenticated user's ID into a request context.
func UserIDToContext(ctx context.Context, userID int64) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}

// AuthMiddleware validates the JWT token in the Authorization header
// AuthMiddleware returns an http.Handler that enforces JWT Bearer authentication and stores the user ID in the request context.
//
// If the Authorization header is missing the middleware responds with HTTP 401 and the body "Missing authorization header".
// If the header is not a Bearer token it responds with HTTP 401 and the body "Invalid authorization header format".
// If token validation fails it responds with HTTP 401 and the body "Invalid or expired token".
// On successful validation the middleware stores the token's UserID in the request context under userIDKey and calls the next handler.
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

// OptionalAuthMiddleware validates the JWT token in the Authorization header if present.
// Unlike AuthMiddleware, it does not require authentication and allows requests to proceed
// without a valid token. If a valid token is present, it stores the user ID in the context.
// If the token is invalid or malformed, the request proceeds without user ID in context.
func OptionalAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader != "" {
			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			if tokenString != authHeader { // Valid Bearer prefix
				claims, err := utils.ValidateJWT(tokenString)
				if err == nil {
					// Token is valid, store user ID in context
					ctx := context.WithValue(r.Context(), userIDKey, claims.UserID)
					r = r.WithContext(ctx)
				}
			}
		}
		next.ServeHTTP(w, r)
	})
}

// UserIDFromContext retrieves the authenticated user ID from the context.
// Returns the ID and true if it exists, 0 and false otherwise.
func UserIDFromContext(ctx context.Context) (int64, bool) {
	id, ok := ctx.Value(userIDKey).(int64)
	return id, ok
}
