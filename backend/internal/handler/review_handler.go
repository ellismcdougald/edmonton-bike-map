package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

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
		// Support either query param (?wayID=) or path segment (/api/reviews/{wayID})
		var wayID int64
		var err error
		wayIDStr := r.URL.Query().Get("wayID")
		if wayIDStr != "" {
			wayID, err = strconv.ParseInt(wayIDStr, 10, 64)
		} else {
			parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
			if len(parts) >= 3 { // [api, reviews, {id}]
				wayID, err = strconv.ParseInt(parts[len(parts)-1], 10, 64)
			} else {
				err = fmt.Errorf("missing wayID")
			}
		}
		if err != nil || wayID == 0 {
			log.Printf("Invalid or missing wayID: %v", err)
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

// HandleDeleteReview deletes the current user's review for a specific way.
// Expects JSON body: { "wayId": number }
// Requires authentication via middleware.
func (h *ReviewHandler) HandleDeleteReview() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type deleteReq struct {
			WayID int64 `json:"wayId"`
		}
		var req deleteReq
		// Body is optional; try path param if not provided
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			// fall through and try path
		}
		if req.WayID == 0 {
			parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
			if len(parts) >= 3 {
				if id, err := strconv.ParseInt(parts[len(parts)-1], 10, 64); err == nil {
					req.WayID = id
				}
			}
		}
		if req.WayID == 0 {
			http.Error(w, "wayId is required", http.StatusBadRequest)
			return
		}

		userID, ok := middleware.UserIDFromContext(r.Context())
		if !ok {
			log.Printf("User ID not found in context")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		if err := h.ReviewService.DeleteUserReviewForWay(userID, req.WayID); err != nil {
			log.Printf("Failed to delete review for user %d way %d: %v", userID, req.WayID, err)
			http.Error(w, "Failed to delete review", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
