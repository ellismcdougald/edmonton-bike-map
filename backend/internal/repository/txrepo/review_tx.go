package txrepo

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/ellismcdougald/edmonton-bike-map/internal/models"
	"github.com/ellismcdougald/edmonton-bike-map/internal/repository"
)

// TxReviewRepository provides review operations within an existing transaction.
type TxReviewRepository struct {
	Tx *sql.Tx
}

// NewTxReviewRepository returns a ReviewRepository that operates using the provided transaction.
func NewTxReviewRepository(tx *sql.Tx) repository.ReviewRepository {
	return &TxReviewRepository{Tx: tx}
}

// CreateReview inserts a review and links it to ways within the transaction.
func (r *TxReviewRepository) CreateReview(review *models.Review) error {
	if review == nil {
		return fmt.Errorf("nil review")
	}
	if len(review.WayIDs) == 0 {
		return fmt.Errorf("review must include at least one way ID")
	}

	insertReview := `
		INSERT INTO reviews (
			user_id,
			rating,
			comment
		) VALUES ($1, $2, $3)
		RETURNING id
	`
	var reviewID int64
	if err := r.Tx.QueryRow(insertReview, review.UserID, review.Rating, review.Comment).Scan(&reviewID); err != nil {
		return err
	}

	linkStmt := `INSERT INTO review_ways (review_id, way_id) VALUES ($1, $2)`
	for _, wayID := range review.WayIDs {
		if _, err := r.Tx.Exec(linkStmt, reviewID, wayID); err != nil {
			return err
		}
	}
	return nil
}

// GetReviews returns reviews associated with a specific way.
func (r *TxReviewRepository) GetReviews(wayID int64) ([]models.Review, error) {
	query := `
		SELECT
			r.user_id,
			r.rating,
			r.comment,
			r.created_at,
			u.username
		FROM reviews r
		JOIN review_ways rw ON rw.review_id = r.id
		LEFT JOIN users u ON r.user_id = u.id
		WHERE rw.way_id = $1;
	`

	rows, err := r.Tx.Query(query, wayID)
	if err != nil {
		return nil, err
	}
	defer func() {
		if cerr := rows.Close(); cerr != nil {
			log.Printf("Error closing rows: %v", cerr)
		}
	}()

	var reviews []models.Review
	for rows.Next() {
		var rev models.Review
		if err := rows.Scan(&rev.UserID, &rev.Rating, &rev.Comment, &rev.CreatedAt, &rev.Username); err != nil {
			return nil, err
		}
		reviews = append(reviews, rev)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if reviews == nil {
		reviews = []models.Review{}
	}

	return reviews, nil
}

// GetAllReviews returns all reviews grouped by way ID via review_ways.
func (r *TxReviewRepository) GetAllReviews() (map[int64][]models.Review, error) {
	query := `
		SELECT
			rw.way_id,
			r.user_id,
			r.rating,
			r.comment,
			r.created_at,
			u.username
		FROM reviews r
		JOIN review_ways rw ON rw.review_id = r.id
		LEFT JOIN users u ON r.user_id = u.id;
	`

	rows, err := r.Tx.Query(query)
	if err != nil {
		return nil, err
	}
	defer func() {
		if cerr := rows.Close(); cerr != nil {
			log.Printf("Error closing rows: %v", cerr)
		}
	}()

	result := make(map[int64][]models.Review)
	for rows.Next() {
		var wayID int64
		var rev models.Review
		if err := rows.Scan(&wayID, &rev.UserID, &rev.Rating, &rev.Comment, &rev.CreatedAt, &rev.Username); err != nil {
			return nil, err
		}
		result[wayID] = append(result[wayID], rev)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

// insertBatch inserts multiple reviews in a single multi-row statement using the transaction.
func (r *TxReviewRepository) insertBatch(reviews []models.Review) error {
	if len(reviews) == 0 {
		return nil
	}

	// Re-implement batch insert via repeated CreateReview calls for clarity.
	for i := range reviews {
		rev := reviews[i]
		if err := r.CreateReview(&rev); err != nil {
			return err
		}
	}
	return nil
}

// InsertBatches inserts reviews in batches of the specified size using the transaction.
func (r *TxReviewRepository) InsertBatches(reviews []models.Review, batchSize int) error {
	if len(reviews) == 0 {
		return nil
	}
	for i := 0; i < len(reviews); i += batchSize {
		end := i + batchSize
		if end > len(reviews) {
			end = len(reviews)
		}
		if err := r.insertBatch(reviews[i:end]); err != nil {
			return err
		}
	}
	return nil
}
