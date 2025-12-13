package sqlrepo

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/ellismcdougald/edmonton-bike-map/internal/models"
	"github.com/ellismcdougald/edmonton-bike-map/internal/repository"
)

// SQLReviewRepository implements ReviewRepository using a SQL database.
type SQLReviewRepository struct {
	DB *sql.DB
}

// NewSQLReviewRepository creates a new SQLReviewRepository.
func NewSQLReviewRepository(db *sql.DB) repository.ReviewRepository {
	return &SQLReviewRepository{DB: db}
}

// CreateReview inserts a new review and its way links (review_ways).
func (r *SQLReviewRepository) CreateReview(review *models.Review) error {
	if review == nil {
		return fmt.Errorf("nil review")
	}
	if len(review.WayIDs) == 0 && review.WayID != 0 {
		review.WayIDs = []int64{review.WayID}
	}
	if len(review.WayIDs) == 0 {
		return fmt.Errorf("review must include at least one way ID")
	}

	// Insert the review and get its ID
	insertReview := `
		INSERT INTO reviews (
			user_id,
			rating,
			comment
		) VALUES ($1, $2, $3)
		RETURNING id
	`
	var reviewID int64
	if err := r.DB.QueryRow(insertReview, review.UserID, review.Rating, review.Comment).Scan(&reviewID); err != nil {
		return err
	}

	// Link the review to its ways
	linkStmt := `INSERT INTO review_ways (review_id, way_id) VALUES ($1, $2)`
	for _, wayID := range review.WayIDs {
		if _, err := r.DB.Exec(linkStmt, reviewID, wayID); err != nil {
			return err
		}
	}
	return nil
}

// GetReviews returns reviews associated with a specific way, including the reviewer's username.
func (r *SQLReviewRepository) GetReviews(wayID int64) ([]models.Review, error) {
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

	rows, err := r.DB.Query(query, wayID)
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

// GetAllReviews retrieves all reviews grouped by way ID via review_ways.
func (r *SQLReviewRepository) GetAllReviews() (map[int64][]models.Review, error) {
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

	rows, err := r.DB.Query(query)
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

// insertBatch inserts multiple reviews and their links; implemented as a simple loop for clarity.
func (r *SQLReviewRepository) insertBatch(reviews []models.Review) error {
	for i := range reviews {
		rev := reviews[i]
		if err := r.CreateReview(&rev); err != nil {
			return err
		}
	}
	return nil
}

// InsertBatches inserts reviews in batches of the specified size.
func (r *SQLReviewRepository) InsertBatches(reviews []models.Review, batchSize int) error {
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
