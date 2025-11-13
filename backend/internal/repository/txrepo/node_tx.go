package txrepo

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/ellismcdougald/edmonton-bike-map/internal/models"
	"github.com/ellismcdougald/edmonton-bike-map/internal/repository"
)

// TxNodeRepository provides node operations within an existing transaction.
type TxNodeRepository struct {
	Tx *sql.Tx
}

// NewTxNodeRepository returns a NodeRepository that operates using the provided transaction.
func NewTxNodeRepository(tx *sql.Tx) repository.NodeRepository {
	return &TxNodeRepository{Tx: tx}
}

// Insert inserts a single node within the transaction.
func (r *TxNodeRepository) Insert(node models.Node) error {
	_, err := r.Tx.Exec(
		"INSERT INTO nodes (id, latitude, longitude) VALUES ($1, $2, $3) ON CONFLICT (id) DO NOTHING",
		node.ID, node.Latitude, node.Longitude,
	)
	return err
}

// insertBatch inserts multiple nodes in a single multi-row statement using the transaction.
func (r *TxNodeRepository) insertBatch(nodes []models.Node) error {
	if len(nodes) == 0 {
		return nil
	}

	query := "INSERT INTO nodes (id, latitude, longitude) VALUES "
	args := []interface{}{}
	for i, n := range nodes {
		if i > 0 {
			query += ", "
		}
		base := i*3 + 1
		query += fmt.Sprintf("($%d, $%d, $%d)", base, base+1, base+2)
		args = append(args, n.ID, n.Latitude, n.Longitude)
	}
	query += " ON CONFLICT (id) DO NOTHING"

	if _, err := r.Tx.Exec(query, args...); err != nil {
		return fmt.Errorf("insert batch nodes: %w", err)
	}
	return nil
}

// InsertBatches inserts nodes in batches of the given size using the transaction.
func (r *TxNodeRepository) InsertBatches(nodes []models.Node, batchSize int) error {
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

// GetNode retrieves a node by ID within the transaction.
func (r *TxNodeRepository) GetNode(id int64) (*models.Node, error) {
	n := &models.Node{}
	err := r.Tx.QueryRow("SELECT id, latitude, longitude FROM nodes WHERE id = $1", id).
		Scan(&n.ID, &n.Latitude, &n.Longitude)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("node %d not found", id)
		}
		return nil, err
	}
	return n, nil
}

// GetAllNodes retrieves all nodes within the transaction.
func (r *TxNodeRepository) GetAllNodes() (map[int64]models.Node, error) {
	rows, err := r.Tx.Query("SELECT id, latitude, longitude FROM nodes")
	if err != nil {
		return nil, err
	}
	defer func() {
		if cerr := rows.Close(); cerr != nil {
			log.Printf("Error closing rows: %v", cerr)
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
