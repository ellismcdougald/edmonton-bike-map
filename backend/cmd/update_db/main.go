package main

import (
	"bytes"
	"database/sql"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/ellismcdougald/edmonton-bike-map/pkg/data"
	"github.com/ellismcdougald/edmonton-bike-map/pkg/model"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

const (
	overpassURL     = "https://overpass-api.de/api/interpreter"
	queryFile       = "cmd/update_db/query.osmql"
	osmJSONPath     = "cmd/update_db/edmonton-osm-data.json"
	nodeBatchSize   = 5000
	reviewBatchSize = 5000
)

func main() {
	loadEnv()
	query := readQueryFile(queryFile)
	fetchAndStoreOSMData(query, osmJSONPath)
	db := connectDB()
	defer closeDB(db)

	osmResp := parseOSMJSON(osmJSONPath)
	processDatabase(db, osmResp)
}

// Load environment variables
func loadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Print("Error loading .env file")
	}
}

// Read OSMQL query file
func readQueryFile(path string) []byte {
	query, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("Could not read query file: %v", err)
	}
	return query
}

// Fetch OSM data from Overpass and store in JSON file
func fetchAndStoreOSMData(query []byte, path string) {
	log.Printf("Fetching data from Overpass API...")
	resp, err := http.Post(overpassURL, "application/x-www-form-urlencoded", bytes.NewReader(query))
	if err != nil {
		log.Fatalf("Could not fetch data: %v", err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			log.Printf("Error closing response body: %v", cerr)
		}
	}()

	tmpPath := path + ".tmp"
	out, err := os.Create(tmpPath)
	if err != nil {
		log.Fatalf("Could not create temp file: %v", err)
	}
	defer func() {
		if cerr := out.Close(); cerr != nil {
			log.Printf("Error closing file: %v", cerr)
		}
	}()

	if _, err := io.Copy(out, resp.Body); err != nil {
		log.Fatalf("Could not write to file: %v", err)
	}

	if err := os.Rename(tmpPath, path); err != nil {
		log.Fatalf("Could not rename file: %v", err)
	}

	log.Printf("OSM data saved to %s", path)
}

// Connect to Postgres database
func connectDB() *sql.DB {
	dbURL := os.Getenv("DATABASE_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Could not connect to DB: %v", err)
	}
	if err := db.Ping(); err != nil {
		log.Fatalf("Could not ping DB: %v", err)
	}
	return db
}

// Close DB with logging
func closeDB(db *sql.DB) {
	if err := db.Close(); err != nil {
		log.Printf("Error closing DB: %v", err)
	}
}

// Parse OSM JSON file
func parseOSMJSON(path string) *data.OSMResponse {
	resp, err := data.ParseOSMJSON(path)
	if err != nil {
		log.Fatalf("Could not parse OSM JSON: %v", err)
	}
	return resp
}

// Process database update
func processDatabase(db *sql.DB, osmResp *data.OSMResponse) {
	tx, err := db.Begin()
	if err != nil {
		log.Fatalf("Failed to begin transaction: %v", err)
	}
	defer func() {
		if rerr := tx.Rollback(); rerr != nil && rerr != sql.ErrTxDone {
			log.Printf("Error rolling back transaction: %v", rerr)
		}
	}()

	// Load existing reviews
	reviewStore := model.TxReviewStore{Tx: tx}
	reviews, err := reviewStore.GetAllReviews()
	if err != nil {
		log.Fatalf("Could not fetch existing reviews: %v", err)
	}
	log.Printf("Loaded %d existing reviews", len(reviews))

	// Clear tables
	clearTables(tx)

	// Prepare nodes and ways
	nodes, ways, wayIDs := extractNodesAndWays(osmResp)

	// Insert nodes
	nodeStore := model.TxNodeStore{Tx: tx}
	if err := nodeStore.InsertBatchChunks(nodes, nodeBatchSize); err != nil {
		log.Fatalf("Could not insert nodes: %v", err)
	}
	log.Printf("Inserted %d nodes", len(nodes))

	// Insert ways
	wayStore := model.TxWayStore{Tx: tx}
	if err := wayStore.InsertBatchDynamic(ways); err != nil {
		log.Fatalf("Could not insert ways: %v", err)
	}
	log.Printf("Inserted %d ways", len(ways))

	// Filter valid reviews
	validReviews := filterValidReviews(reviews, wayIDs)
	if err := reviewStore.InsertBatchChunks(validReviews, reviewBatchSize); err != nil {
		log.Fatalf("Could not insert reviews: %v", err)
	}
	log.Printf("Inserted %d reviews", len(validReviews))

	// Commit transaction
	if err := tx.Commit(); err != nil {
		log.Fatalf("Failed to commit transaction: %v", err)
	}
}

// Clear relevant tables
func clearTables(tx *sql.Tx) {
	_, err := tx.Exec("TRUNCATE reviews, way_nodes, ways, nodes RESTART IDENTITY CASCADE;")
	if err != nil {
		log.Fatalf("Could not clear tables: %v", err)
	}
	log.Print("Cleared tables successfully")
}

// Extract nodes and ways from OSM response
func extractNodesAndWays(osmResp *data.OSMResponse) ([]model.DBNode, []model.DBWay, map[int64]struct{}) {
	nodes := []model.DBNode{}
	ways := []model.DBWay{}
	wayIDs := make(map[int64]struct{})

	for _, el := range osmResp.Elements {
		switch el.Type {
		case "node":
			nodes = append(nodes, model.DBNode{ID: el.ID, Latitude: el.Lat, Longitude: el.Lon})
		case "way":
			ways = append(ways, model.DBWay{ID: el.ID, Tags: el.Tags, NodeIDs: el.Nodes})
			wayIDs[el.ID] = struct{}{}
		}
	}

	return nodes, ways, wayIDs
}

// Filter reviews to only include valid way IDs
func filterValidReviews(reviews []model.Review, wayIDs map[int64]struct{}) []model.Review {
	valid := []model.Review{}
	for _, r := range reviews {
		if _, ok := wayIDs[r.WayID]; ok {
			valid = append(valid, r)
		} else {
			log.Printf("Skipping review for unknown way ID %d", r.WayID)
		}
	}
	return valid
}
