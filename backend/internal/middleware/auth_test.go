package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ellismcdougald/edmonton-bike-map/internal/utils"
	"github.com/stretchr/testify/require"
)

func TestOptionalAuthMiddleware_ValidToken(t *testing.T) {
	utils.SetJWTKey([]byte("test-jwt-key"))
	token, err := utils.GenerateJWT("testuser", 42)
	require.NoError(t, err)

	var capturedUserID int64
	var foundUserID bool

	handler := OptionalAuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedUserID, foundUserID = UserIDFromContext(r.Context())
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	require.True(t, foundUserID)
	require.Equal(t, int64(42), capturedUserID)
}

func TestOptionalAuthMiddleware_NoToken(t *testing.T) {
	utils.SetJWTKey([]byte("test-jwt-key"))

	var foundUserID bool

	handler := OptionalAuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, foundUserID = UserIDFromContext(r.Context())
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	require.False(t, foundUserID, "user ID should not be present in context without token")
}

func TestOptionalAuthMiddleware_InvalidToken(t *testing.T) {
	utils.SetJWTKey([]byte("test-jwt-key"))

	var foundUserID bool

	handler := OptionalAuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, foundUserID = UserIDFromContext(r.Context())
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer invalid-token-xyz")
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	require.False(t, foundUserID, "user ID should not be present in context with invalid token")
}

func TestOptionalAuthMiddleware_MalformedAuthHeader(t *testing.T) {
	utils.SetJWTKey([]byte("test-jwt-key"))

	var foundUserID bool

	handler := OptionalAuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, foundUserID = UserIDFromContext(r.Context())
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "NotBearer sometoken")
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	require.False(t, foundUserID, "user ID should not be present in context with malformed header")
}

func TestOptionalAuthMiddleware_AllowsRequestToProceed(t *testing.T) {
	utils.SetJWTKey([]byte("test-jwt-key"))

	handlerCalled := false

	handler := OptionalAuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlerCalled = true
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	require.True(t, handlerCalled, "next handler should be called even without authentication")
	require.Equal(t, http.StatusOK, rr.Code)
}

func TestUserIDFromContext_Missing(t *testing.T) {
	ctx := context.Background()
	userID, found := UserIDFromContext(ctx)

	require.False(t, found)
	require.Equal(t, int64(0), userID)
}

func TestUserIDFromContext_Present(t *testing.T) {
	ctx := UserIDToContext(context.Background(), 123)
	userID, found := UserIDFromContext(ctx)

	require.True(t, found)
	require.Equal(t, int64(123), userID)
}

func TestUserIDToContext(t *testing.T) {
	ctx := context.Background()
	newCtx := UserIDToContext(ctx, 99)

	require.NotEqual(t, ctx, newCtx)

	userID, found := UserIDFromContext(newCtx)
	require.True(t, found)
	require.Equal(t, int64(99), userID)
}
