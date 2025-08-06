package handlers

import (
	"encoding/json"
	"net/http"

	"golang.org/x/crypto/bcrypt"

	"github.com/ellismcdougald/edmonton-bike-map/pkg/model"
	"github.com/ellismcdougald/edmonton-bike-map/pkg/utils"
)

func (h* RealHandlers) HandleLogin() http.HandlerFunc {
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

func (h* RealHandlers) HandleSignup() http.HandlerFunc {
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

/*
// HandleLogin handles HTTP requests to log in to the application.
// The request body includes the given username and password
// The username and password is validated against existing usernames and (hashed) passwords
// Responds with a generated JWT token if login is successful
// Returns HTTP 401 if username or password is rejected
var HandleLogin = func(writer http.ResponseWriter, request *http.Request, db *sql.DB) {
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

	user, err := model.GetUser(db, creds.Username)
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

// HandleSignUp handles HTTP requests to sign up for the application
// Request body includes proposed username and password
// If username does not already exist, then the new user is created
// Returns HTTP 201 on success.
var HandleSignUp = func(writer http.ResponseWriter, request *http.Request, db *sql.DB) {
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

	exists, err := model.UsernameExists(db, creds.Username)
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

	err = user.Create(db)
	if err != nil {
		http.Error(writer, "Error creating user", http.StatusInternalServerError)
		return
	}

	writer.WriteHeader(http.StatusCreated)
}
*/
