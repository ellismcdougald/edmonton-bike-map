package models

import "time"

type Review struct {
	// Deprecated: WayID will be removed after multi-way migration is complete.
	// Repositories will treat WayID as a single-element WayIDs when provided.
	WayID     int64
	WayIDs    []int64
	UserID    int64
	Rating    int
	Comment   string
	CreatedAt time.Time
	Username  string
}
