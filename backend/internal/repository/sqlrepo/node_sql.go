package sqlrepo

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/ellismcdougald/edmonton-bike-map/internal/models"
	"github.com/ellismcdougald/edmonton-bike-map/internal/repository"
)

// SQLNodeRepository provides methods to interact with the nodes table in the database.
type SQLNodeRepository struct {
	DB *sql.DB
}

// NewSQLNodeRepository creates a new instance of SQLNodeRepository.
func NewSQLNodeRepository(db *sql.DB) repository.NodeRepository {
	return &SQLNodeRepository{DB: db}
}

// Insert adds a single node to the database.
func (r *SQLNodeRepository) Insert(node models.Node) error {
	_, err := r.DB.Exec(
		"INSERT INTO nodes (id, latitude, longitude) VALUES ($1, $2, $3) ON CONFLICT (id) DO NOTHING",
		node.ID, node.Latitude, node.Longitude,
	)
	return err
}

// insertBatch adds multiple nodes to the database in a single query.
// helper used by InsertBatches.
func (r *SQLNodeRepository) insertBatch(nodes []models.Node) error {
	if len(nodes) == 0 {
		return nil
	}

	// Start a transaction
	tx, err := r.DB.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err := tx.Rollback(); err != nil && err != sql.ErrTxDone {
			log.Printf("tx rollback error: %v", err)
		}
	}()

	// Build the multi-row insert statement
	query := "INSERT INTO nodes (id, latitude, longitude) VALUES "
	args := []interface{}{}
	for i, n := range nodes {
		if i > 0 {
			query += ", "
		}
		// $1, $2, $3, ... for each node
		base := i*3 + 1
		query += fmt.Sprintf("($%d, $%d, $%d)", base, base+1, base+2)
		args = append(args, n.ID, n.Latitude, n.Longitude)
	}
	query += " ON CONFLICT (id) DO NOTHING"

	// Execute the batch insert
	if _, err := tx.Exec(query, args...); err != nil {
		return err
	}

	return tx.Commit()
}

// InsertBatches inserts nodes in batches of the specified size.
func (r *SQLNodeRepository) InsertBatches(nodes []models.Node, batchSize int) error {
	if len(nodes) == 0 {
		return nil
	}

	for i := 0; i < len(nodes); i += batchSize {
		end := i + batchSize
		if end > len(nodes) {
			end = len(nodes)
		}
		if err := r.insertBatch(nodes[i:end]); err != nil {
			return err
		}
	}

	return nil
}

// GetNode retrieves a node by its ID.
// Returns a pointer to the Node and an error if not found.
func (r *SQLNodeRepository) GetNode(id int64) (*models.Node, error) {
	n := &models.Node{}
	err := r.DB.QueryRow("SELECT id, latitude, longitude FROM nodes WHERE id = $1", id).
		Scan(&n.ID, &n.Latitude, &n.Longitude)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("node not found")
		}
		return nil, err
	}
	return n, nil
}

// GetAllNodes retrieves all nodes from the database.
// Returns a map of node ID to Node and an error if the query fails.
func (r *SQLNodeRepository) GetAllNodes() (map[int64]models.Node, error) {
	query := `
		SELECT 
			id,
			latitude,
			longitude
		FROM nodes;
	`

	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("Error closing rows: %v", err)
		}
	}()

	nodes := make(map[int64]models.Node)
	for rows.Next() {
		var node models.Node
		if err := rows.Scan(&node.ID, &node.Latitude, &node.Longitude); err != nil {
			return nil, err
		}
		nodes[node.ID] = node
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return nodes, nil
}
