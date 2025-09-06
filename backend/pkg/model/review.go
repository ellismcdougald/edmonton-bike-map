package model

import (
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
