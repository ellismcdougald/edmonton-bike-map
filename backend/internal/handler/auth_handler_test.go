package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"golang.org/x/crypto/bcrypt"

	"github.com/ellismcdougald/edmonton-bike-map/internal/models"
	"github.com/ellismcdougald/edmonton-bike-map/internal/service"
	"github.com/ellismcdougald/edmonton-bike-map/internal/token"
	"github.com/stretchr/testify/require"
)

// mockUserRepo implements repository.UserRepository for tests.
type mockUserRepo struct {
	user              *models.User
	getErr            error
	usernameExists    bool
	usernameErr       error
	createErr         error
	updatePasswordErr error
}

func (m *mockUserRepo) GetByUsername(username string) (*models.User, error) {
	return m.user, m.getErr
}

func (m *mockUserRepo) GetByID(id int64) (*models.User, error) {
	return m.user, m.getErr
}

func (m *mockUserRepo) Create(user *models.User) error {
	return m.createErr
}

func (m *mockUserRepo) UsernameExists(username string) (bool, error) {
	return m.usernameExists, m.usernameErr
}

func (m *mockUserRepo) UpdatePassword(userID int64, hashedPassword string) error {
	return m.updatePasswordErr
}

// fake provider
type fakeProvider struct {
	token string
	err   error
}

func (f *fakeProvider) Generate(username string, userID int64) (string, error) {
	return f.token, f.err
}

// compile-time check that fakeProvider implements token.TokenProvider
var _ token.TokenProvider = (*fakeProvider)(nil)

func TestHandleLogin_Success(t *testing.T) {
	pw := []byte("secret")
	hash, err := bcrypt.GenerateFromPassword(pw, bcrypt.DefaultCost)
	require.NoError(t, err)

	repo := &mockUserRepo{user: &models.User{ID: 1, Username: "alice", Password: string(hash)}}
	svc := service.NewUserService(repo)
	fp := &fakeProvider{token: "fixed-token"}

	h := NewAuthHandler(svc, fp)

	body, _ := json.Marshal(models.Credentials{Username: "alice", Password: "secret"})
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
	rr := httptest.NewRecorder()

	h.HandleLogin().ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	require.Equal(t, "application/json", rr.Header().Get("Content-Type"))

	var resp map[string]string
	require.NoError(t, json.NewDecoder(rr.Body).Decode(&resp))
	require.Equal(t, "fixed-token", resp["token"])
}

func TestHandleLogin_WrongPassword(t *testing.T) {
	hash, err := bcrypt.GenerateFromPassword([]byte("other"), bcrypt.DefaultCost)
	require.NoError(t, err)

	repo := &mockUserRepo{user: &models.User{ID: 1, Username: "alice", Password: string(hash)}}
	svc := service.NewUserService(repo)
	fp := &fakeProvider{token: "x"}
	h := NewAuthHandler(svc, fp)

	body, _ := json.Marshal(models.Credentials{Username: "alice", Password: "secret"})
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
	rr := httptest.NewRecorder()

	h.HandleLogin().ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestHandleLogin_BadMethod(t *testing.T) {
	repo := &mockUserRepo{}
	svc := service.NewUserService(repo)
	fp := &fakeProvider{token: "x"}
	h := NewAuthHandler(svc, fp)

	req := httptest.NewRequest(http.MethodGet, "/login", nil)
	rr := httptest.NewRecorder()

	h.HandleLogin().ServeHTTP(rr, req)

	require.Equal(t, http.StatusMethodNotAllowed, rr.Code)
}

func TestHandleSignup_Success(t *testing.T) {
	repo := &mockUserRepo{usernameExists: false}
	svc := service.NewUserService(repo)
	fp := &fakeProvider{token: "jwt"}
	h := NewAuthHandler(svc, fp)

	body, _ := json.Marshal(models.Credentials{Username: "bob", Password: "pw"})
	req := httptest.NewRequest(http.MethodPost, "/signup", bytes.NewReader(body))
	rr := httptest.NewRecorder()

	h.HandleSignup().ServeHTTP(rr, req)

	require.Equal(t, http.StatusCreated, rr.Code)
}

func TestHandleSignup_UsernameExists(t *testing.T) {
	repo := &mockUserRepo{usernameExists: true}
	svc := service.NewUserService(repo)
	fp := &fakeProvider{token: "jwt"}
	h := NewAuthHandler(svc, fp)

	body, _ := json.Marshal(models.Credentials{Username: "bob", Password: "pw"})
	req := httptest.NewRequest(http.MethodPost, "/signup", bytes.NewReader(body))
	rr := httptest.NewRecorder()

	h.HandleSignup().ServeHTTP(rr, req)

	require.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestHandleSignup_BadRequest(t *testing.T) {
	repo := &mockUserRepo{}
	svc := service.NewUserService(repo)
	fp := &fakeProvider{token: "jwt"}
	h := NewAuthHandler(svc, fp)

	req := httptest.NewRequest(http.MethodPost, "/signup", bytes.NewReader([]byte("not json")))
	rr := httptest.NewRecorder()

	h.HandleSignup().ServeHTTP(rr, req)

	require.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestHandleSignup_CreateError(t *testing.T) {
	repo := &mockUserRepo{usernameExists: false, createErr: errors.New("db gone")}
	svc := service.NewUserService(repo)
	fp := &fakeProvider{token: "jwt"}
	h := NewAuthHandler(svc, fp)

	body, _ := json.Marshal(models.Credentials{Username: "bob", Password: "pw"})
	req := httptest.NewRequest(http.MethodPost, "/signup", bytes.NewReader(body))
	rr := httptest.NewRecorder()

	h.HandleSignup().ServeHTTP(rr, req)

	require.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestHandleChangePassword_Success(t *testing.T) {
	// Set JWT key for token validation
	testSecret := []byte("test-secret")
	jwtProvider := token.NewJWTProvider(testSecret)

	pw := []byte("oldpassword")
	hash, err := bcrypt.GenerateFromPassword(pw, bcrypt.DefaultCost)
	require.NoError(t, err)

	repo := &mockUserRepo{user: &models.User{ID: 123, Username: "alice", Password: string(hash)}}
	svc := service.NewUserService(repo)
	h := NewAuthHandler(svc, jwtProvider)

	// Generate a valid token
	validToken, err := jwtProvider.Generate("alice", 123)
	require.NoError(t, err)

	body, _ := json.Marshal(models.ChangePasswordRequest{
		CurrentPassword: "oldpassword",
		NewPassword:     "newpassword",
	})
	req := httptest.NewRequest(http.MethodPost, "/change-password", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+validToken)
	rr := httptest.NewRecorder()

	h.HandleChangePassword().ServeHTTP(rr, req)

	require.Equal(t, http.StatusNoContent, rr.Code)
}

func TestHandleChangePassword_WrongCurrentPassword(t *testing.T) {
	testSecret := []byte("test-secret")
	jwtProvider := token.NewJWTProvider(testSecret)

	pw := []byte("oldpassword")
	hash, err := bcrypt.GenerateFromPassword(pw, bcrypt.DefaultCost)
	require.NoError(t, err)

	repo := &mockUserRepo{user: &models.User{ID: 123, Username: "alice", Password: string(hash)}}
	svc := service.NewUserService(repo)
	h := NewAuthHandler(svc, jwtProvider)

	validToken, err := jwtProvider.Generate("alice", 123)
	require.NoError(t, err)

	body, _ := json.Marshal(models.ChangePasswordRequest{
		CurrentPassword: "wrongpassword",
		NewPassword:     "newpassword",
	})
	req := httptest.NewRequest(http.MethodPost, "/change-password", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+validToken)
	rr := httptest.NewRecorder()

	h.HandleChangePassword().ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestHandleChangePassword_MissingToken(t *testing.T) {
	testSecret := []byte("test-secret")
	jwtProvider := token.NewJWTProvider(testSecret)

	repo := &mockUserRepo{}
	svc := service.NewUserService(repo)
	h := NewAuthHandler(svc, jwtProvider)

	body, _ := json.Marshal(models.ChangePasswordRequest{
		CurrentPassword: "oldpassword",
		NewPassword:     "newpassword",
	})
	req := httptest.NewRequest(http.MethodPost, "/change-password", bytes.NewReader(body))
	rr := httptest.NewRecorder()

	h.HandleChangePassword().ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnauthorized, rr.Code)
}
