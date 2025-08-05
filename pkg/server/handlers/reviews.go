package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/ellismcdougald/edmonton-bike-map/pkg/model"
)

// HandleGetReviews handles HTTP requests to get all reviews associated with a certain way
// Query paramemters: wayID (int64)
// Responds with a list of Review objects
var HandleGetReviews = func(writer http.ResponseWriter, request *http.Request, db *sql.DB) {
	writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
	writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	writer.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")

	query := request.URL.Query()
	wayID, err := strconv.ParseInt(query.Get("wayID"), 10, 64)
	if err != nil {
		log.Printf("Error extracting parameter wayID from query %v: %v", query, err)
		http.Error(writer, "Invalid wayID"+err.Error(), http.StatusBadRequest)
		return
	}

	reviews, err := model.GetReviews(db, wayID)
	if err != nil {
		log.Printf("Error fetching reviews from database for wayID %d: %v", wayID, err)
		http.Error(writer, "Could not fetch reviews"+err.Error(), http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(writer).Encode(reviews)
	if err != nil {
		log.Printf("Error encoding json")
		http.Error(writer, "Could not encode json", http.StatusInternalServerError)
	}
}

// HandlePostReview handles HTTP requests to create a review for a certain way
// Request body contains review details
// The review is added to the database
// Returns HTTP 400 if review details are invalid
var HandlePostReview = func(writer http.ResponseWriter, request *http.Request, db *sql.DB) {
	writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
	writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")

	if request.Method == http.MethodOptions {
		writer.WriteHeader(http.StatusOK)
		return
	}

	if request.Method != http.MethodPost {
		http.Error(writer, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var review model.Review
	if err := json.NewDecoder(request.Body).Decode(&review); err != nil {
		log.Printf("Error decoding review JSON: %v", err)
		http.Error(writer, "Bad Request", http.StatusBadRequest)
		return
	}

	if review.Rating < 1 || review.Rating > 10 {
		log.Printf("Invalid review: Rating must be between 1 and 10 inclusive")
		http.Error(writer, "Invalid review: Rating must be between 1 and 10 inclusive", http.StatusBadRequest)
		return
	}

	err := review.Create(db)
	if err != nil {
		http.Error(writer, "Database error", http.StatusInternalServerError)
		return
	}
}
