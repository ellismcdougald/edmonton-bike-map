package models

import "time"

type Review struct {
	ID        int64
	WayIDs    []int64
	UserID    int64
	Rating    int
	Comment   string
	CreatedAt time.Time
	Username  string
}
