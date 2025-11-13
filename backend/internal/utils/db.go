package utils

import (
	"database/sql"
	"log"
)

func Connect(dbURL string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

func Close(db *sql.DB) {
	if err := db.Close(); err != nil {
		log.Printf("Error closing DB: %v", err)
	}
}
