package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/ellismcdougald/edmonton-bike-map/pkg/model"
)

// HandleGetReviews returns an HTTP handler that fetches all reviews for a given way.
// It expects a query parameter "wayID" specifying the way's ID.
// Responds with a JSON array of reviews on success,
// HTTP 400 Bad Request if the wayID is missing or invalid,
// and HTTP 500 Internal Server Error if fetching reviews fails.
func (h *RealHandlers) HandleGetReviews() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		wayIDStr := query.Get("wayID")
		wayID, err := strconv.ParseInt(wayIDStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid wayID", http.StatusBadRequest)
			return
		}

		reviews, err := h.ReviewService.GetReviews(wayID)
		if err != nil {
			http.Error(w, "Failed to fetch reviews", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(reviews); err != nil {
			http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
		}
	}
}

// HandlePostReview returns an HTTP handler that accepts a new review submission via POST.
// Supports CORS preflight OPTIONS requests.
// Expects JSON body with rating (1-10), way ID, user ID, and comment.
// Responds with HTTP 201 Created on success,
// HTTP 400 Bad Request for invalid input,
// and HTTP 500 Internal Server Error if saving the review fails.
func (h *RealHandlers) HandlePostReview() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		var review model.Review
		if err := json.NewDecoder(r.Body).Decode(&review); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if review.Rating < 1 || review.Rating > 10 {
			http.Error(w, "Rating must be between 1 and 10", http.StatusBadRequest)
			return
		}

		err := h.ReviewService.CreateReview(&review)
		if err != nil {
			http.Error(w, "Failed to save review", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
	}
}
