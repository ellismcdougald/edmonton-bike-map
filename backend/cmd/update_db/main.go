package main

import (
	"log"
	"os"

	"github.com/ellismcdougald/edmonton-bike-map/internal/updater"
	"github.com/ellismcdougald/edmonton-bike-map/internal/utils"
)

const (
	queryFile = "cmd/update_db/query.osmql"
)

func main() {
	utils.LoadEnv()

	query := readQueryFile(queryFile)

	data, err := updater.FetchOverpassData(query)
	if err != nil {
		log.Fatalf("Error fetching Overpass data: %v", err)
	}

	osmResp, err := updater.ParseOSMBytes(data)
	if err != nil {
		log.Fatalf("Error parsing OSM data: %v", err)
	}

	db, err := utils.Connect(os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	defer utils.Close(db)

	err = updater.UpdateDatabase(db, osmResp)
	if err != nil {
		log.Fatalf("Error updating database: %v", err)
	}
}

func readQueryFile(path string) []byte {
	query, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("Could not read query file: %v", err)
	}
	return query
}
