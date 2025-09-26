package model

import (
	"database/sql"
	"fmt"
	"log"
)

// TxReviewStore provides methods to insert and retrieve reviews within an existing transaction.
type TxReviewStore struct {
	Tx *sql.Tx
}

// GetAllReviews retrieves all reviews, including the username of each reviewer.
func (s *TxReviewStore) GetAllReviews() ([]Review, error) {
	query := `
		SELECT
			r.way_id,
			r.user_id,
			r.rating,
			r.comment,
			r.created_at,
			u.username
		FROM reviews r
		LEFT JOIN users u ON r.user_id = u.id;
	`

	rows, err := s.Tx.Query(query)
	if err != nil {
		return nil, err
	}
	defer func() {
		if cerr := rows.Close(); cerr != nil {
			log.Printf("Error closing rows: %v", cerr)
		}
	}()

	var reviews []Review
	for rows.Next() {
		var review Review
		if err := rows.Scan(&review.WayID, &review.UserID, &review.Rating, &review.Comment, &review.CreatedAt, &review.Username); err != nil {
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

func (s *TxReviewStore) GetReviews(wayID int64) ([]Review, error) {
	query := `
		SELECT
			r.way_id,
			r.user_id,        -- Add this line
			r.rating,
			r.comment,
			r.created_at,
			u.username
		FROM reviews r
		LEFT JOIN users u ON r.user_id = u.id
		WHERE r.way_id = $1;
	`
	rows, err := s.Tx.Query(query, wayID)
	if err != nil {
		return nil, err
	}
	defer func() {
		if cerr := rows.Close(); cerr != nil {
			log.Printf("Error closing rows: %v", cerr)
		}
	}()

	var reviews []Review
	for rows.Next() {
		var review Review
		err = rows.Scan(&review.WayID, &review.UserID, &review.Rating, &review.Comment, &review.CreatedAt, &review.Username)  // Add &review.UserID
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

// InsertReview inserts a single review using the current transaction.
func (s *TxReviewStore) InsertReview(review *Review) error {
	query := `
	INSERT INTO reviews (
		way_id,
		user_id,
		rating,
		comment
	)
	VALUES ($1, $2, $3, $4)
	`
	_, err := s.Tx.Exec(query, review.WayID, review.UserID, review.Rating, review.Comment)
	return err
}

// InsertBatch inserts multiple reviews in a single multi-row statement within the transaction.
func (s *TxReviewStore) InsertBatch(reviews []Review) error {
	if len(reviews) == 0 {
		return nil
	}

	query := "INSERT INTO reviews (way_id, user_id, rating, comment) VALUES "
	args := []interface{}{}
	for i, r := range reviews {
		if i > 0 {
			query += ", "
		}
		base := i*4 + 1
		query += fmt.Sprintf("($%d, $%d, $%d, $%d)", base, base+1, base+2, base+3)
		args = append(args, r.WayID, r.UserID, r.Rating, r.Comment)
	}

	_, err := s.Tx.Exec(query, args...)
	return err
}

// InsertBatchChunks inserts reviews in smaller batches to avoid exceeding parameter limits.
// Each batch is inserted using InsertBatch.
func (s *TxReviewStore) InsertBatchChunks(reviews []Review, batchSize int) error {
	if len(reviews) == 0 {
		return nil
	}
	if batchSize <= 0 {
		batchSize = 1000 // default batch size
	}

	for i := 0; i < len(reviews); i += batchSize {
		end := i + batchSize
		if end > len(reviews) {
			end = len(reviews)
		}
		batch := reviews[i:end]
		if err := s.InsertBatch(batch); err != nil {
			return fmt.Errorf("failed to insert review batch %d-%d: %w", i, end, err)
		}
	}

	return nil
}
