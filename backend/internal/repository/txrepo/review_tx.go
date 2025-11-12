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

// CreateReview inserts a review using the transaction.
func (r *TxReviewRepository) CreateReview(review *models.Review) error {
    query := `
    INSERT INTO reviews (
        way_id,
        user_id,
        rating,
        comment
    ) VALUES ($1, $2, $3, $4)
    `
    _, err := r.Tx.Exec(query, review.WayID, review.UserID, review.Rating, review.Comment)
    return err
}

// GetReviews returns reviews for a specific way.
func (r *TxReviewRepository) GetReviews(wayID int64) ([]models.Review, error) {
    query := `
        SELECT
            r.way_id,
            r.user_id,
            r.rating,
            r.comment,
            r.created_at,
            u.username
        FROM reviews r
        LEFT JOIN users u ON r.user_id = u.id
        WHERE r.way_id = $1;
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
        if err := rows.Scan(&rev.WayID, &rev.UserID, &rev.Rating, &rev.Comment, &rev.CreatedAt, &rev.Username); err != nil {
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

// GetAllReviews returns all reviews grouped by way ID.
func (r *TxReviewRepository) GetAllReviews() (map[int64][]models.Review, error) {
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
        var rev models.Review
        if err := rows.Scan(&rev.WayID, &rev.UserID, &rev.Rating, &rev.Comment, &rev.CreatedAt, &rev.Username); err != nil {
            return nil, err
        }
        result[rev.WayID] = append(result[rev.WayID], rev)
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

    query := "INSERT INTO reviews (way_id, user_id, rating, comment) VALUES "
    args := []interface{}{}
    for i, rev := range reviews {
        if i > 0 {
            query += ", "
        }
        base := i*4 + 1
        query += fmt.Sprintf("($%d, $%d, $%d, $%d)", base, base+1, base+2, base+3)
        args = append(args, rev.WayID, rev.UserID, rev.Rating, rev.Comment)
    }

    if _, err := r.Tx.Exec(query, args...); err != nil {
        return err
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
