package model

import (
	"database/sql"
	"log"
	"time"
)

// Review represents a review for a way with rating, comment, and creation timestamp.
type Review struct {
	WayID     int64     `json:"wayId"`
	UserID    int64     `json:"userId"`
	Rating    int       `json:"rating"`
	Comment   string    `json:"comment"`
	CreatedAt time.Time `json:"createdAt"`
	Username  string    `json:"username"`
}

// Create inserts the Review into the database.
func (review *Review) Create(db *sql.DB) error {
	query := `
	INSERT INTO reviews (
		way_id,
		user_id,
		rating,
		comment
	)
	VALUES (
		$1,
		$2,
		$3,
		$4
	)
	`
	_, err := db.Exec(query, review.WayID, review.UserID, review.Rating, review.Comment)
	return err
}

// GetReviews retrieves all reviews for a specific way by wayID from the database.
// Returns a slice of Review and any error encountered.
func GetReviews(db *sql.DB, wayID int64) ([]Review, error) {
	query := `
		SELECT
  		r.way_id,
  		r.rating,
  		r.comment,
  		r.created_at,
  		u.username
		FROM reviews r
		LEFT JOIN users u ON r.user_id = u.id
		WHERE r.way_id = $1;
	`

	rows, err := db.Query(query, wayID)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("Error closing rows: %v", err)
		}
	}()

	var reviews []Review
	for rows.Next() {
		var review Review
		err = rows.Scan(&review.WayID, &review.Rating, &review.Comment, &review.CreatedAt, &review.Username)
		if err != nil {
			return nil, err
		}
		reviews = append(reviews, review)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return reviews, nil
}
