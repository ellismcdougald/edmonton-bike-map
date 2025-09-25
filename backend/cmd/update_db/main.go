package main

/*
TODOS:
- Need to clear the nodes, ways, and way_nodes tables
- Need to implement transactions for database inserts
- Need to make sure other tables are protected from cascading deletes
*/

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
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Print("Error loading .env file")
	}

	overpassURL := "https://overpass-api.de/api/interpreter"
	osmJSONPath := "cmd/update_db/edmonton-osm-data.json"

	// Fetch most recent osm data
	queryFile := "cmd/update_db/query.osmql"
	query, err := os.ReadFile(queryFile)
	if err != nil {
		log.Fatalf("Could not read query file: %v", err)
	}

	log.Printf("Fetching query from Overpass")
	resp, err := http.Post(overpassURL, "application/x-www-form-urlencoded", bytes.NewReader(query))
	if err != nil {
		log.Fatalf("Could not get query from Overpass: %v", err)
	}
	defer resp.Body.Close()
	log.Printf("Done fetching data from overpass")

	tmpPath := osmJSONPath + ".tmp"
	out, err := os.Create(tmpPath)
	if err != nil {
		log.Fatalf("Could not create file to store result: %v", err)
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		log.Printf("Could not store result in file: %v", err)
	}

	os.Rename(tmpPath, osmJSONPath)
	log.Printf("Data written to file")

	// Connect to database
	dbURL := os.Getenv("DATABASE_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Could not connect to db: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("Error closing db: %v", err)
		}
	}()
	if err := db.Ping(); err != nil {
		log.Fatalf("Could not ping db: %v", err)
	}

	osmResp, err := data.ParseOSMJSON(osmJSONPath)
	if err != nil {
		log.Fatalf("Could not parse osm data: %v", err)
	}

	// Begin updates
	tx, err := db.Begin()
	if err != nil {
		log.Fatalf("Failed to begin transaction: %v", err)
	}
	defer tx.Rollback()

	reviewStore := model.TxReviewStore{Tx: tx}
	reviews, err := reviewStore.GetAllReviews()
	if err != nil {
		log.Fatalf("Could not get reviews from database: %v", err)
	}
	log.Printf("Saved existing reviews successfully")

	// Clear reviews, nodes, ways, and way_nodes
	_, err = tx.Exec("TRUNCATE reviews, way_nodes, ways, nodes RESTART IDENTITY CASCADE;")
	if err != nil {
		log.Fatalf("Could not clear tables: %v", err)
	}
	log.Printf("Cleared tables successfully")
	
	nodes := []model.DBNode{}
	ways := []model.DBWay{}
	wayIDs := make(map[int64]struct{})
	for _, el := range osmResp.Elements {
		if el.Type == "node" {
			n := model.DBNode{
				ID:        el.ID,
				Latitude:  el.Lat,
				Longitude: el.Lon,
			}
			nodes = append(nodes, n)
		} else if el.Type == "way" {
			w := model.DBWay{
				ID:      el.ID,
				Tags:    el.Tags,
				NodeIDs: el.Nodes,
			}
			ways = append(ways, w)
			wayIDs[el.ID] = struct{}{}
		}
	}
	nodeStore := model.TxNodeStore{Tx: tx}
	err = nodeStore.InsertBatchChunks(nodes, 5000)
	if err != nil {
		log.Fatalf("Could not insert nodes: %v", err)
	}
	log.Printf("Inserted nodes succesfully")
	wayStore := model.TxWayStore{Tx: tx}
	err = wayStore.InsertBatchDynamic(ways)
	if err != nil {
		log.Fatalf("Could not insert ways: %v", err)
	}
	log.Printf("Inserted ways succesfully")

	// Insert reviews
	validReviews := []model.Review{}
	for _, r := range reviews {
		if _, ok := wayIDs[r.WayID]; ok {
			validReviews = append(validReviews, r)
		} else {
			log.Printf("Skipping review for unknown way id %d", r.WayID)
		}
	}
	err = reviewStore.InsertBatchChunks(validReviews, 5000)
	if err != nil {
		log.Fatalf("Could not insert reviews: %v", err)
	}
	log.Printf("Inserted reviews succesfully")

	if err = tx.Commit(); err != nil {
    log.Fatalf("Failed to commit transaction: %v", err)
	}












	


}