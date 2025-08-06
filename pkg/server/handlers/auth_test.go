package handlers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ellismcdougald/edmonton-bike-map/pkg/model"
	"github.com/ellismcdougald/edmonton-bike-map/pkg/server/handlers"
	"github.com/ellismcdougald/edmonton-bike-map/pkg/utils"
	"golang.org/x/crypto/bcrypt"
)

// Define MockUserStore:
type MockUserStore struct {
	GetUserFunc func(username string) (*model.User, error)
	UsernameExistsFunc func(username string) (bool, error)
	CreateUserFunc func(user *model.User) error
}

func (m *MockUserStore) GetUser(username string) (*model.User, error) {
	if m.GetUserFunc != nil {
		return m.GetUserFunc(username)
	}
	return nil, errors.New("GetUserFunc not implemented")
}

func (m *MockUserStore) UsernameExists(username string) (bool, error) {
	if m.UsernameExistsFunc != nil {
		return m.UsernameExistsFunc(username)
	}
	return false, errors.New("UsernameExistsFunc not implemented")
}

func (m *MockUserStore) CreateUser(user *model.User) error {
	if m.CreateUserFunc != nil {
		return m.CreateUserFunc(user)
	}
	return errors.New("CreateUserFunc not implemented")
}

// HandleLogin tests:
func TestHandleLogin(t *testing.T) {
	password := "password123"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	utils.SetJWTKey([]byte("mysecretkey"))

	tests := []struct {
		name           string
		method         string
		requestBody    map[string]string
		mockGetUser    func(username string) (*model.User, error)
		wantStatusCode int
		wantToken      bool
	}{
		{
			name:   "Successful login",
			method: http.MethodPost,
			requestBody: map[string]string{
				"username": "testuser",
				"password": password,
			},
			mockGetUser: func(username string) (*model.User, error) {
				return &model.User{
					Username: username,
					Password: string(hashedPassword),
					ID:       42,
				}, nil
			},
			wantStatusCode: http.StatusOK,
			wantToken:      true,
		},
		{
			name:           "Invalid method",
			method:         http.MethodGet,
			requestBody:    nil,
			mockGetUser:    nil,
			wantStatusCode: http.StatusMethodNotAllowed,
			wantToken:      false,
		},
		{
			name:   "Invalid JSON",
			method: http.MethodPost,
			requestBody: nil,
			mockGetUser:    nil,
			wantStatusCode: http.StatusBadRequest,
			wantToken:      false,
		},
		{
			name:   "User not found",
			method: http.MethodPost,
			requestBody: map[string]string{
				"username": "nonexistent",
				"password": "any",
			},
			mockGetUser: func(username string) (*model.User, error) {
				return nil, errors.New("user not found")
			},
			wantStatusCode: http.StatusUnauthorized,
			wantToken:      false,
		},
		{
			name:   "Wrong password",
			method: http.MethodPost,
			requestBody: map[string]string{
				"username": "testuser",
				"password": "wrongpassword",
			},
			mockGetUser: func(username string) (*model.User, error) {
				return &model.User{
					Username: username,
					Password: string(hashedPassword),
					ID:       42,
				}, nil
			},
			wantStatusCode: http.StatusUnauthorized,
			wantToken:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var req *http.Request
			if tt.requestBody != nil {
				bodyBytes, err := json.Marshal(tt.requestBody)
				if err != nil {
					t.Fatalf("Failed to marshal requestBody: %v", err)
				}
				req = httptest.NewRequest(tt.method, "/api/login", bytes.NewReader(bodyBytes))
			} else if tt.name == "Invalid JSON" {
				// purposely give invalid JSON by passing bad reader
				req = httptest.NewRequest(tt.method, "/api/login", bytes.NewReader([]byte("{invalid-json")))
			} else {
				req = httptest.NewRequest(tt.method, "/api/login", nil)
			}

			rec := httptest.NewRecorder()

			mockStore := &MockUserStore{
				GetUserFunc: tt.mockGetUser,
			}

			realHandlers := handlers.RealHandlers{
				UserService: mockStore,
			}

			handler := realHandlers.HandleLogin()
			handler(rec, req)

			resp := rec.Result()
			defer resp.Body.Close()

			if resp.StatusCode != tt.wantStatusCode {
				t.Errorf("Status code = %d; want %d", resp.StatusCode, tt.wantStatusCode)
			}

			if tt.wantToken {
				var respBody map[string]string
				if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
					t.Fatalf("Failed to decode response body: %v", err)
				}

				token, ok := respBody["token"]
				if !ok || token == "" {
					t.Errorf("Expected token in response but none found")
				}
			}
		})
	}
}

func TestHandleSignUp(t *testing.T) {
	utils.SetJWTKey([]byte("mysecretkey"))

	tests := []struct {
		name             string
		method           string
		requestBody      map[string]string
		mockExists       func(username string) (bool, error)
		mockCreateUser   func(user *model.User) error
		wantStatusCode   int
	}{
		{
			name:   "Successful signup",
			method: http.MethodPost,
			requestBody: map[string]string{
				"username": "newuser",
				"password": "securepass",
			},
			mockExists: func(username string) (bool, error) {
				return false, nil
			},
			mockCreateUser: func(user *model.User) error {
				return nil
			},
			wantStatusCode: http.StatusCreated,
		},
		{
			name:   "Username already exists",
			method: http.MethodPost,
			requestBody: map[string]string{
				"username": "existing",
				"password": "somepass",
			},
			mockExists: func(username string) (bool, error) {
				return true, nil
			},
			mockCreateUser: nil,
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name:           "Invalid method",
			method:         http.MethodGet,
			requestBody:    nil,
			wantStatusCode: http.StatusMethodNotAllowed,
		},
		{
			name:           "Invalid JSON",
			method:         http.MethodPost,
			requestBody:    nil, // trigger invalid JSON case
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name:   "Missing username or password",
			method: http.MethodPost,
			requestBody: map[string]string{
				"username": "",
				"password": "",
			},
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name:   "Database error on create",
			method: http.MethodPost,
			requestBody: map[string]string{
				"username": "erroruser",
				"password": "securepass",
			},
			mockExists: func(username string) (bool, error) {
				return false, nil
			},
			mockCreateUser: func(user *model.User) error {
				return errors.New("db create failed")
			},
			wantStatusCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var req *http.Request
			if tt.name == "Invalid JSON" {
				req = httptest.NewRequest(tt.method, "/api/signup", bytes.NewReader([]byte("{bad json")))
			} else if tt.requestBody != nil {
				bodyBytes, _ := json.Marshal(tt.requestBody)
				req = httptest.NewRequest(tt.method, "/api/signup", bytes.NewReader(bodyBytes))
			} else {
				req = httptest.NewRequest(tt.method, "/api/signup", nil)
			}

			rec := httptest.NewRecorder()

			mockStore := &MockUserStore{
				UsernameExistsFunc: tt.mockExists,
				CreateUserFunc:     tt.mockCreateUser,
			}

			realHandlers := handlers.RealHandlers{
				UserService: mockStore,
			}

			handler := realHandlers.HandleSignup()
			handler(rec, req)

			resp := rec.Result()
			defer resp.Body.Close()

			if resp.StatusCode != tt.wantStatusCode {
				t.Errorf("Status = %d; want %d", resp.StatusCode, tt.wantStatusCode)
			}
		})
	}
}
