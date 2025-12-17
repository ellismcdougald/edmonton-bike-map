package updater

import (
	"database/sql"
	"log"
	"time"

	"github.com/ellismcdougald/edmonton-bike-map/internal/models"
	"github.com/ellismcdougald/edmonton-bike-map/internal/repository/txrepo"
	"github.com/lib/pq"
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
	defer func() {
		if rerr := tx.Rollback(); rerr != nil && rerr != sql.ErrTxDone {
			log.Printf("transaction rollback error: %v", rerr)
		}
	}()

	reviewRepo := txrepo.NewTxReviewRepository(tx)
	// Snapshot existing reviews with their linked way IDs before clearing tables
	reviews, err := snapshotExistingReviews(tx)
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

	// Prune review way mappings to only those that still exist after repopulating ways
	validReviews := pruneReviewsToExistingWays(reviews, wayIDs)
	if err := reviewRepo.InsertBatches(validReviews, reviewBatchSize); err != nil {
		return err
	}
	log.Printf("Inserted %d reviews", len(validReviews))

	// Ensure reviews.id sequence is set past the max id to avoid conflicts on future inserts
	if err := adjustReviewIDSequence(tx); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	log.Print("Database update committed successfully")
	return nil
}

// adjustReviewIDSequence sets the reviews.id sequence to the current max(id)
// to ensure subsequent inserts without explicit IDs do not conflict.
func adjustReviewIDSequence(tx *sql.Tx) error {
	_, err := tx.Exec(`
		SELECT setval(
			pg_get_serial_sequence('reviews', 'id'),
			COALESCE((SELECT MAX(id) FROM reviews), 0)
		)
	`)
	return err
}

// snapshotExistingReviews captures all reviews and their associated way IDs in the current DB state.
// It groups by review ID to ensure each review's way links are preserved distinctly.
func snapshotExistingReviews(tx *sql.Tx) ([]models.Review, error) {
	query := `
		SELECT
			r.id,
			r.user_id,
			r.rating,
			r.comment,
			r.created_at,
			u.username,
			COALESCE(array_agg(rw.way_id ORDER BY rw.way_id)
					 FILTER (WHERE rw.way_id IS NOT NULL), '{}') AS way_ids
		FROM reviews r
		LEFT JOIN review_ways rw ON rw.review_id = r.id
		LEFT JOIN users u ON r.user_id = u.id
		GROUP BY r.id, r.user_id, r.rating, r.comment, r.created_at, u.username;
	`

	rows, err := tx.Query(query)
	if err != nil {
		return nil, err
	}
	defer func() {
		if cerr := rows.Close(); cerr != nil {
			log.Printf("Error closing rows: %v", cerr)
		}
	}()

	var out []models.Review
	for rows.Next() {
		var (
			reviewID  int64
			userID    int64
			rating    int
			comment   string
			createdAt time.Time
			username  string
			wayIDs    pq.Int64Array
		)
		if err := rows.Scan(&reviewID, &userID, &rating, &comment, &createdAt, &username, &wayIDs); err != nil {
			return nil, err
		}
		out = append(out, models.Review{
			ID:        reviewID,
			WayIDs:    []int64(wayIDs),
			UserID:    userID,
			Rating:    rating,
			Comment:   comment,
			CreatedAt: createdAt,
			Username:  username,
		})
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return out, nil
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

// pruneReviewsToExistingWays removes any way IDs from each review that no longer exist;
// reviews that end up with zero way IDs are dropped.
func pruneReviewsToExistingWays(reviews []models.Review, wayIDs map[int64]struct{}) []models.Review {
	valid := make([]models.Review, 0, len(reviews))
	for _, r := range reviews {
		pruned := make([]int64, 0, len(r.WayIDs))
		for _, wid := range r.WayIDs {
			if _, ok := wayIDs[wid]; ok {
				pruned = append(pruned, wid)
			} else {
				log.Printf("Skipping non-existent way %d for review by user %d", wid, r.UserID)
			}
		}
		if len(pruned) == 0 {
			// No valid way links remain; drop this review
			continue
		}
		r.WayIDs = pruned
		valid = append(valid, r)
	}
	return valid
}
