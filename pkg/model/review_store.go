package model

import (
	"database/sql"
	"log"
)

// DBReviewStore provides methods to interact with the reviews table in the database.
type DBReviewStore struct {
	DB *sql.DB
}

// Create inserts the Review into the database.
func (s *DBReviewStore) CreateReview(review *Review) error {
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
	_, err := s.DB.Exec(query, review.WayID, review.UserID, review.Rating, review.Comment)
	return err
}

// GetReviews retrieves all reviews for a specific way by wayID from the database.
// Returns a slice of Review and any error encountered.
func (s *DBReviewStore) GetReviews(wayID int64) ([]Review, error) {
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

	rows, err := s.DB.Query(query, wayID)
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

	if reviews == nil {
		reviews = []Review{}
	}

	return reviews, nil
}
