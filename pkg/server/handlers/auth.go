package handlers

import (
	"encoding/json"
	"net/http"

	"golang.org/x/crypto/bcrypt"

	"github.com/ellismcdougald/edmonton-bike-map/pkg/model"
	"github.com/ellismcdougald/edmonton-bike-map/pkg/utils"
)

// HandleLogin returns an HTTP handler that processes user login requests.
// It expects POST requests with JSON body containing username and password.
// On successful authentication, it returns a JSON Web Token (JWT).
// Responds with HTTP 401 Unauthorized if credentials are invalid,
// and HTTP 405 Method Not Allowed if the method is not POST.
func (h *RealHandlers) HandleLogin() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodPost {
			http.Error(writer, "Only POST allowed", http.StatusMethodNotAllowed)
			return
		}

		var creds model.Credentials
		err := json.NewDecoder(request.Body).Decode(&creds)
		if err != nil {
			http.Error(writer, "Invalid request body", http.StatusBadRequest)
			return
		}

		user, err := h.UserService.GetUser(creds.Username)
		if err != nil {
			http.Error(writer, "Invalid username or password", http.StatusUnauthorized)
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password)); err != nil {
			http.Error(writer, "Invalid username or password", http.StatusUnauthorized)
			return
		}

		token, err := utils.GenerateJWT(user.Username, user.ID)
		if err != nil {
			http.Error(writer, "Could not generate token", http.StatusInternalServerError)
			return
		}

		writer.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(writer).Encode(map[string]string{
			"token": token,
		})
		if err != nil {
			http.Error(writer, "Could not encode JWT token", http.StatusInternalServerError)
			return
		}
	}
}

// HandleSignup returns an HTTP handler that processes user signup requests.
// It expects POST requests with JSON body containing username and password.
// Returns HTTP 201 Created on successful user creation,
// HTTP 400 Bad Request for missing fields or duplicate usernames,
// and HTTP 405 Method Not Allowed if the method is not POST.
func (h *RealHandlers) HandleSignup() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodPost {
			http.Error(writer, "Only POST allowed", http.StatusMethodNotAllowed)
			return
		}

		var creds model.Credentials
		err := json.NewDecoder(request.Body).Decode(&creds)
		if err != nil {
			http.Error(writer, "Invalid request body", http.StatusBadRequest)
			return
		}

		if creds.Username == "" || creds.Password == "" {
			http.Error(writer, "Username and password required", http.StatusBadRequest)
		}

		exists, err := h.UserService.UsernameExists(creds.Username)
		if err != nil {
			http.Error(writer, "Database error", http.StatusInternalServerError)
			return
		}
		if exists {
			http.Error(writer, "Username already taken", http.StatusBadRequest)
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(creds.Password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(writer, "Error hashing password", http.StatusInternalServerError)
			return
		}

		user := &model.User{
			Username: creds.Username,
			Password: string(hashedPassword),
		}

		err = h.UserService.CreateUser(user)
		if err != nil {
			http.Error(writer, "Error creating user", http.StatusInternalServerError)
			return
		}

		writer.WriteHeader(http.StatusCreated)
	}
}