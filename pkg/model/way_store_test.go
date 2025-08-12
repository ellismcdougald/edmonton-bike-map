package model_test

import (
	"encoding/json"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/ellismcdougald/edmonton-bike-map/pkg/model"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestDBWayStore_Insert(t *testing.T) {
	ways := []model.DBWay{
		{
			ID:      1,
			Tags:    map[string]string{"highway": "residential"},
			NodeIDs: []int64{10, 20, 30},
		},
		{
			ID:      2,
			Tags:    map[string]string{"highway": "primary"},
			NodeIDs: []int64{100, 200},
		},
	}

	for _, way := range ways {
		t.Run("Insert_Success", func(t *testing.T) {
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer db.Close()

			store := &model.DBWayStore{DB: db}

			tagsJSON, _ := json.Marshal(way.Tags)

			mock.ExpectBegin()

			mock.ExpectExec(regexp.QuoteMeta("INSERT INTO ways (id, tags) VALUES ($1, $2) ON CONFLICT (id) DO NOTHING")).
				WithArgs(way.ID, tagsJSON).
				WillReturnResult(sqlmock.NewResult(1, 1))

			for seq, nodeID := range way.NodeIDs {
				mock.ExpectExec(regexp.QuoteMeta("INSERT INTO way_nodes (way_id, node_id, sequence_id) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING")).
					WithArgs(way.ID, nodeID, seq).
					WillReturnResult(sqlmock.NewResult(1, 1))
			}

			mock.ExpectCommit()

			err = store.Insert(way)
			assert.NoError(t, err)
			assert.NoError(t, mock.ExpectationsWereMet())
		})

		t.Run("Insert_WaysInsertError", func(t *testing.T) {
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer db.Close()

			store := &model.DBWayStore{DB: db}

			tagsJSON, _ := json.Marshal(way.Tags)

			mock.ExpectBegin()

			mock.ExpectExec(regexp.QuoteMeta("INSERT INTO ways (id, tags) VALUES ($1, $2) ON CONFLICT (id) DO NOTHING")).
				WithArgs(way.ID, tagsJSON).
				WillReturnError(errors.New("insert ways error"))

			mock.ExpectRollback()

			err = store.Insert(way)
			assert.Error(t, err)
			assert.NoError(t, mock.ExpectationsWereMet())
		})

		t.Run("Insert_WayNodesInsertError", func(t *testing.T) {
			if len(way.NodeIDs) == 0 {
				t.Skip("No nodes to test way_nodes insert error")
			}
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer db.Close()

			store := &model.DBWayStore{DB: db}

			tagsJSON, _ := json.Marshal(way.Tags)

			mock.ExpectBegin()

			mock.ExpectExec(regexp.QuoteMeta("INSERT INTO ways (id, tags) VALUES ($1, $2) ON CONFLICT (id) DO NOTHING")).
				WithArgs(way.ID, tagsJSON).
				WillReturnResult(sqlmock.NewResult(1, 1))

			// Simulate success for first few nodes, error for the last
			for seq, nodeID := range way.NodeIDs {
				if seq < len(way.NodeIDs)-1 {
					mock.ExpectExec(regexp.QuoteMeta("INSERT INTO way_nodes (way_id, node_id, sequence_id) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING")).
						WithArgs(way.ID, nodeID, seq).
						WillReturnResult(sqlmock.NewResult(1, 1))
				} else {
					mock.ExpectExec(regexp.QuoteMeta("INSERT INTO way_nodes (way_id, node_id, sequence_id) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING")).
						WithArgs(way.ID, nodeID, seq).
						WillReturnError(errors.New("insert way_nodes error"))
					mock.ExpectRollback()
				}
			}

			err = store.Insert(way)
			assert.Error(t, err)
			assert.NoError(t, mock.ExpectationsWereMet())
		})

		t.Run("Insert_BeginError", func(t *testing.T) {
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer db.Close()

			store := &model.DBWayStore{DB: db}

			mock.ExpectBegin().WillReturnError(errors.New("begin transaction error"))

			err = store.Insert(way)
			assert.Error(t, err)
			assert.NoError(t, mock.ExpectationsWereMet())
		})

		t.Run("Insert_CommitError", func(t *testing.T) {
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer db.Close()

			store := &model.DBWayStore{DB: db}

			tagsJSON, _ := json.Marshal(way.Tags)

			mock.ExpectBegin()

			mock.ExpectExec(regexp.QuoteMeta("INSERT INTO ways (id, tags) VALUES ($1, $2) ON CONFLICT (id) DO NOTHING")).
				WithArgs(way.ID, tagsJSON).
				WillReturnResult(sqlmock.NewResult(1, 1))

			for seq, nodeID := range way.NodeIDs {
				mock.ExpectExec(regexp.QuoteMeta("INSERT INTO way_nodes (way_id, node_id, sequence_id) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING")).
					WithArgs(way.ID, nodeID, seq).
					WillReturnResult(sqlmock.NewResult(1, 1))
			}

			mock.ExpectCommit().WillReturnError(errors.New("commit error"))

			err = store.Insert(way)
			assert.Error(t, err)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestDBWayStore_GetAllWays(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	store := &model.DBWayStore{DB: db}

	ways := []model.DBWay{
		{
			ID:      1,
			Tags:    map[string]string{"highway": "residential"},
			NodeIDs: []int64{10, 20, 30},
		},
		{
			ID:      2,
			Tags:    map[string]string{"highway": "primary"},
			NodeIDs: []int64{100, 200},
		},
	}

	// Prepare rows to be returned by mock query
	rows := sqlmock.NewRows([]string{"id", "tags", "node_ids"})
	for _, w := range ways {
		tagsJSON, _ := json.Marshal(w.Tags)
		rows.AddRow(w.ID, tagsJSON, pq.Array(w.NodeIDs))
	}

	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT
			w.id,
			w.tags,
			ARRAY_AGG(wn.node_id ORDER BY wn.sequence_id) AS node_ids
		FROM ways w
		JOIN way_nodes wn ON wn.way_id = w.id
		GROUP BY w.id, w.tags;
	`)).
		WillReturnRows(rows)

	gotWays, err := store.GetAllWays()
	assert.NoError(t, err)
	assert.Equal(t, ways, gotWays)
	assert.NoError(t, mock.ExpectationsWereMet())

	// Test query error
	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT
			w.id,
			w.tags,
			ARRAY_AGG(wn.node_id ORDER BY wn.sequence_id) AS node_ids
		FROM ways w
		JOIN way_nodes wn ON wn.way_id = w.id
		GROUP BY w.id, w.tags;
	`)).
		WillReturnError(errors.New("query failed"))

	_, err = store.GetAllWays()
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
