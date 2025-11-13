package main

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"

	"github.com/ellismcdougald/edmonton-bike-map/internal/models"
	"github.com/ellismcdougald/edmonton-bike-map/internal/repository/sqlrepo"
	"github.com/ellismcdougald/edmonton-bike-map/internal/updater"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Print("Error loading .env file")
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatalf("No DATABASE_URL given")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Could not connect to db: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("Error closing db: %v", err)
		}
	}()

	filePath, err := filepath.Abs(filepath.Join(".", "osm_bike_data.json"))
	if err != nil {
		log.Fatalf("cannot get absolute path: %v", err)
	}
	data, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Could not read osm data file: %v", err)
	}
	resp, err := updater.ParseOSMBytes(data)
	if err != nil {
		log.Fatalf("Could not parse osm data: %v", err)
	}

	nodeRepo := sqlrepo.NewSQLNodeRepository(db)
	wayRepo := sqlrepo.NewSQLWayRepository(db)

	var nodes []models.Node
	var ways []models.Way

	for _, el := range resp.Elements {
		switch el.Type {
		case "node":
			nodes = append(nodes, models.Node{
				ID:        el.ID,
				Latitude:  el.Lat,
				Longitude: el.Lon,
			})
		case "way":
			ways = append(ways, models.Way{
				ID:      el.ID,
				Tags:    el.Tags,
				NodeIDs: el.Nodes,
			})
		}
	}

	log.Printf("Inserting nodes")
	// Batch insert all nodes
	if err := nodeRepo.InsertBatches(nodes, 1000); err != nil {
		log.Fatalf("failed to insert nodes in batches: %v", err)
	}

	log.Printf("Inserting ways")
	// Batch insert all ways
	if err := wayRepo.InsertBatches(ways, 500); err != nil {
		log.Fatalf("failed to insert ways in batches: %v", err)
	}
	log.Printf("done populating")
}
