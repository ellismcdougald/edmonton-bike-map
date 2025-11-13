package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	"github.com/ellismcdougald/edmonton-bike-map/internal/models"
	"github.com/ellismcdougald/edmonton-bike-map/internal/repository/sqlrepo"
	"github.com/ellismcdougald/edmonton-bike-map/internal/updater"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("POSTGRES_DB")
	port := os.Getenv("POSTGRES_PORT")

	connectionStr := fmt.Sprintf("postgres://%s:%s@localhost:%s/%s?sslmode=disable", user, password, port, dbname)

	db, err := sql.Open("postgres", connectionStr)
	if err != nil {
		log.Fatalf("Could not connect to db: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("Error closing db: %v", err)
		}
	}()

	data, err := os.ReadFile("osm_bike_data.json")
	if err != nil {
		log.Fatalf("Could not read osm data file: %v", err)
	}
	resp, err := updater.ParseOSMBytes(data)
	if err != nil {
		log.Fatalf("Could not parse osm data: %v", err)
	}

	nodeRepo := sqlrepo.NewSQLNodeRepository(db)
	for _, el := range resp.Elements {
		if el.Type == "node" {
			n := models.Node {
				ID: 			el.ID,
				Latitude: 	el.Lat,
				Longitude: el.Lon,
			}
			err = nodeRepo.Insert(n)
			if err != nil {
				log.Printf("Warning: inserting node %d failed with error: %v", el.ID, err)
			}
		}
	}

	wayRepo := sqlrepo.NewSQLWayRepository(db)
	for _, el := range resp.Elements {
		if el.Type == "way" {
			w := models.Way{
				ID:      el.ID,
				Tags:    el.Tags,
				NodeIDs: el.Nodes,
			}
			err = wayRepo.Insert(w)
			if err != nil {
				log.Printf("Warning: inserting way %d failed with error: %v", el.ID, err)
			}
		}
	}
}
