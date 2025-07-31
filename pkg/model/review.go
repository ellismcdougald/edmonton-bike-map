package model

type Review struct {
	WayID			 int64	`json:"wayId"`
	Rating     int    `json:"rating"`
	ReviewText string `json:"reviewText"`
}
