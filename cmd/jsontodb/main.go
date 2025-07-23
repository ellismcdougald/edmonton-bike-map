package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq" // <-- Postgres driver import

	"github.com/ellismcdougald/edmonton-bike-map/pkg/data"
	"github.com/ellismcdougald/edmonton-bike-map/pkg/model"
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
	defer db.Close()

	resp, err := data.ParseOSMJSON("osm_bike_data.json")
	if err != nil {
		log.Fatalf("Could not parse osm data: %v", err)
	}

	for _, el := range resp.Elements {
		if el.Type == "node" {
			n := model.DBNode{
				ID: el.ID, 
				Latitude: el.Lat, 
				Longitude: el.Lon,
			}
			err = n.Insert(db)
			if err != nil {
				log.Printf("Warning: inserting node %d failed with error: %v", el.ID, err)
			}
		}
	}

	for _, el := range resp.Elements {
		if el.Type == "way" {
			w := model.DBWay{
				ID: el.ID,
				Tags: el.Tags,
				NodeIDs: el.Nodes,
			}
			err = w.Insert(db)
			if err != nil {
				log.Printf("Warning: inserting way %d failed with error: %v", el.ID, err)
			}
		}
	}
}