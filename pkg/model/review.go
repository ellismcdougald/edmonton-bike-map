package model

import (
	"database/sql"
	"time"
)

type Review struct {
	WayID			 int64	`json:"wayId"`
	Rating     int    `json:"rating"`
	ReviewText string `json:"reviewText"`
	CreatedAt	 time.Time `json:"createdAt"`
}

func GetReviews(db *sql.DB, wayID int64) ([]Review, error) {
	query := `
		SELECT
			r.way_id,
			r.rating,
			r.comment
			r.created_at
		FROM reviews r
		WHERE r.way_id = ?
	`

	rows, err := db.Query(query, wayID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reviews []Review
	for rows.Next() {
		var review Review
		
		err = rows.Scan(&review.WayID, &review.Rating, &review.ReviewText, &review.CreatedAt)
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
