package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/ellismcdougald/edmonton-bike-map/internal/service"
	"github.com/ellismcdougald/edmonton-bike-map/internal/utils"
)

// UserHandler handles user-specific HTTP requests (settings etc.).
type UserHandler struct {
	UserService *service.UserService
}

// NewUserHandler creates a new UserHandler.
func NewUserHandler(userSvc *service.UserService) *UserHandler {
	return &UserHandler{UserService: userSvc}
}

// HandleGetSettings returns an HTTP handler that responds with the current user's settings.
func (h *UserHandler) HandleGetSettings() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Only GET allowed", http.StatusMethodNotAllowed)
			return
		}

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
			log.Printf("user_handler: token validation failed: %v", err)
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		user, err := h.UserService.GetUserByID(claims.UserID)
		if err != nil || user == nil {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		resp := map[string]interface{}{
			"username":     user.Username,
			"cyclingSpeed": user.CyclingSpeed,
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	}
}

// HandleUpdateCyclingSpeed returns an HTTP handler that updates the user's cycling speed.
func (h *UserHandler) HandleUpdateCyclingSpeed() http.HandlerFunc {
	type reqBody struct {
		CyclingSpeed int `json:"cyclingSpeed"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
			return
		}

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing authorization header", http.StatusUnauthorized)
			return
		}

		tokenString := authHeader[len("Bearer "):]
		claims, err := utils.ValidateJWT(tokenString)
		if err != nil {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		var body reqBody
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if body.CyclingSpeed <= 0 || body.CyclingSpeed > 80 {
			http.Error(w, "Cycling speed must be between 1 and 80 km/h", http.StatusBadRequest)
			return
		}

		if err := h.UserService.UpdateCyclingSpeed(claims.UserID, body.CyclingSpeed); err != nil {
			log.Printf("failed to update cycling speed for user %d: %v", claims.UserID, err)
			http.Error(w, "Failed to update settings", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
