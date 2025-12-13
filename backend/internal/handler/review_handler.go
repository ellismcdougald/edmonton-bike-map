package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/ellismcdougald/edmonton-bike-map/internal/middleware"
	"github.com/ellismcdougald/edmonton-bike-map/internal/models"
	"github.com/ellismcdougald/edmonton-bike-map/internal/service"
)

type ReviewHandler struct {
	ReviewService *service.ReviewService
}

func NewReviewHandler(reviewService *service.ReviewService) *ReviewHandler {
	return &ReviewHandler{
		ReviewService: reviewService,
	}
}

// HandleGetReviews returns an HTTP handler that fetches all reviews for a given way.
// It expects a query parameter "wayID" specifying the way's ID.
// Responds with a JSON array of reviews on success,
// HTTP 400 Bad Request if the wayID is missing or invalid,
// and HTTP 500 Internal Server Error if fetching reviews fails.
func (h *ReviewHandler) HandleGetReviews() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		wayIDStr := r.URL.Query().Get("wayID")
		wayID, err := strconv.ParseInt(wayIDStr, 10, 64)
		if err != nil {
			log.Printf("Invalid wayID: %v", err)
			http.Error(w, "Invalid wayID", http.StatusBadRequest)
			return
		}

		reviews, err := h.ReviewService.GetReviewsForWay(wayID)
		if err != nil {
			log.Printf("Failed to fetch reviews for way %d: %v", wayID, err)
			http.Error(w, "Failed to fetch reviews", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(reviews); err != nil {
			log.Printf("Failed to encode reviews JSON: %v", err)
			http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
		}
	}
}

// HandlePostReview returns an HTTP handler that accepts a new review submission via POST.
// Expects JSON body with rating (1-10), way ID, user ID, and comment.
// Responds with HTTP 201 Created on success,
// HTTP 400 Bad Request for invalid input,
// and HTTP 500 Internal Server Error if saving the review fails.
func (h *ReviewHandler) HandlePostReview() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Accept either legacy single wayId or new wayIds array
		type postReviewRequest struct {
			WayIDs  []int64 `json:"wayIds"`
			WayID   *int64  `json:"wayId"`
			Rating  int     `json:"rating"`
			Comment string  `json:"comment"`
		}
		var req postReviewRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Printf("Invalid request body: %v", err)
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Get user ID from context (set by AuthMiddleware)
		userID, ok := middleware.UserIDFromContext(r.Context())
		if !ok {
			log.Printf("User ID not found in context")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		// Build models.Review from request, mapping legacy field if needed
		var review models.Review
		review.UserID = userID
		review.Rating = req.Rating
		review.Comment = req.Comment
		if len(req.WayIDs) > 0 {
			review.WayIDs = req.WayIDs
		} else if req.WayID != nil {
			review.WayIDs = []int64{*req.WayID}
		}

		if review.Rating < 1 || review.Rating > 10 {
			log.Printf("Invalid rating: %d", review.Rating)
			http.Error(w, "Rating must be between 1 and 10", http.StatusBadRequest)
			return
		}

		if len(review.WayIDs) == 0 {
			log.Printf("Missing wayIds")
			http.Error(w, "At least one wayId is required", http.StatusBadRequest)
			return
		}

		if err := h.ReviewService.AddReview(review); err != nil {
			log.Printf("Failed to save review: %v", err)
			http.Error(w, "Failed to save review", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
	}
}
