package updater

import (
	"database/sql"
	"log"

	"github.com/ellismcdougald/edmonton-bike-map/internal/models"
	"github.com/ellismcdougald/edmonton-bike-map/internal/repository/txrepo"
)

const (
	nodeBatchSize   = 100
	reviewBatchSize = 1000
)

func UpdateDatabase(db *sql.DB, osmResp *OSMResponse) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	reviewRepo := txrepo.NewTxReviewRepository(tx)
	reviews, err := reviewRepo.GetAllReviews()
	if err != nil {
		return err
	}

	clearTables(tx)

	nodes, ways, wayIDs := extractNodesAndWays(osmResp)

	nodeRepo := txrepo.NewTxNodeRepository(tx)
	if err := nodeRepo.InsertBatches(nodes, nodeBatchSize); err != nil {
		return err
	}
	log.Printf("Inserted %d nodes", len(nodes))

	wayRepo := txrepo.NewTxWayRepository(tx)
	if err := wayRepo.InsertBatches(ways, nodeBatchSize); err != nil {
		return err
	}
	log.Printf("Inserted %d ways", len(ways))

	validReviews := filterValidReviews(reviews, wayIDs)
	if err := reviewRepo.InsertBatches(validReviews, reviewBatchSize); err != nil {
		return err
	}
	log.Printf("Inserted %d reviews", len(validReviews))

	if err := tx.Commit(); err != nil {
		return err
	}
	log.Print("Database update committed successfully")
	return nil
}

func clearTables(tx *sql.Tx) {
	_, err := tx.Exec("TRUNCATE reviews, way_nodes, ways, nodes RESTART IDENTITY CASCADE;")
	if err != nil {
		log.Fatalf("Could not clear tables: %v", err)
	}
	log.Print("Cleared tables successfully")
}

func extractNodesAndWays(osmResp *OSMResponse) ([]models.Node, []models.Way, map[int64]struct{}) {
	nodes := []models.Node{}
	ways := []models.Way{}
	wayIDSet := make(map[int64]struct{})

	for _, el := range osmResp.Elements {
		switch el.Type {
		case "node":
			nodes = append(nodes, models.Node{ID: el.ID, Latitude: el.Lat, Longitude: el.Lon})
		case "way":
			ways = append(ways, models.Way{ID: el.ID, Tags: el.Tags, NodeIDs: el.Nodes})
			wayIDSet[el.ID] = struct{}{}
		}
	}
	return nodes, ways, wayIDSet
}

func filterValidReviews(reviewMap map[int64][]models.Review, wayIDs map[int64]struct{}) []models.Review {
    valid := []models.Review{}
    for wayID, revs := range reviewMap {
        if _, exists := wayIDs[wayID]; exists {
            valid = append(valid, revs...)
        } else {
			for _, r := range revs {
				log.Printf("Skipping review for non-existent way %d: %+v", wayID, r)
			}
        }
    }
    return valid
}
