package txrepo

import (
	"database/sql"
	"encoding/json"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/ellismcdougald/edmonton-bike-map/internal/models"
	"github.com/ellismcdougald/edmonton-bike-map/internal/repository"
	"github.com/stretchr/testify/require"
)

func TestTxWayRepository_Insert(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer func() { mock.ExpectClose(); _ = db.Close() }()

	mock.ExpectBegin()
	tx, err := db.Begin()
	require.NoError(t, err)

	repo := NewTxWayRepository(tx)

	tags := map[string]string{"highway": "cycleway"}
	tagsJSON, _ := json.Marshal(tags)

	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO ways (id, tags) VALUES ($1, $2) ON CONFLICT (id) DO NOTHING")).
		WithArgs(int64(1), tagsJSON).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO way_nodes (way_id, node_id, sequence_id) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING")).
		WithArgs(int64(1), int64(10), 0).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO way_nodes (way_id, node_id, sequence_id) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING")).
		WithArgs(int64(1), int64(20), 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Insert(models.Way{ID: 1, Tags: tags, NodeIDs: []int64{10, 20}})
	require.NoError(t, err)

	mock.ExpectCommit()
	require.NoError(t, tx.Commit())
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestTxWayRepository_GetWay(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer func() { mock.ExpectClose(); _ = db.Close() }()

	mock.ExpectBegin()
	tx, err := db.Begin()
	require.NoError(t, err)

	repo := NewTxWayRepository(tx)

	tags := map[string]string{"highway": "cycleway"}
	tagsJSON, _ := json.Marshal(tags)

	rows := sqlmock.NewRows([]string{"id", "tags"}).AddRow(int64(1), tagsJSON)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, tags FROM ways WHERE id = $1")).WithArgs(int64(1)).WillReturnRows(rows)

	nodeRows := sqlmock.NewRows([]string{"node_id"}).AddRow(int64(10)).AddRow(int64(20))
	mock.ExpectQuery(regexp.QuoteMeta("SELECT node_id FROM way_nodes WHERE way_id = $1 ORDER BY sequence_id")).WithArgs(int64(1)).WillReturnRows(nodeRows)

	way, err := repo.GetWay(1)
	require.NoError(t, err)
	require.Equal(t, int64(1), way.ID)
	require.Equal(t, tags, way.Tags)
	require.Equal(t, []int64{10, 20}, way.NodeIDs)

	mock.ExpectCommit()
	require.NoError(t, tx.Commit())
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestTxWayRepository_GetAllWays(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer func() { mock.ExpectClose(); _ = db.Close() }()

	mock.ExpectBegin()
	tx, err := db.Begin()
	require.NoError(t, err)

	repo := NewTxWayRepository(tx)

	tags1, _ := json.Marshal(map[string]string{"a": "1"})
	tags2, _ := json.Marshal(map[string]string{"b": "2"})
	rows := sqlmock.NewRows([]string{"id", "tags", "node_ids"}).
		AddRow(int64(1), tags1, "{1,2}").
		AddRow(int64(2), tags2, "{3,4}")
	mock.ExpectQuery(regexp.QuoteMeta(`
        SELECT
            w.id,
            w.tags,
            ARRAY_AGG(wn.node_id ORDER BY wn.sequence_id) AS node_ids
        FROM ways w
        JOIN way_nodes wn ON wn.way_id = w.id
        GROUP BY w.id, w.tags;
    `)).WillReturnRows(rows)

	ways, err := repo.GetAllWays()
	require.NoError(t, err)
	require.Len(t, ways, 2)

	mock.ExpectCommit()
	require.NoError(t, tx.Commit())
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestTxWayRepository_InsertBatches(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer func() { mock.ExpectClose(); _ = db.Close() }()

	mock.ExpectBegin()
	tx, err := db.Begin()
	require.NoError(t, err)

	repo := NewTxWayRepository(tx)

	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO ways (id, tags) VALUES ($1, $2), ($3, $4) ON CONFLICT (id) DO NOTHING")).
		WithArgs(int64(3), sqlmock.AnyArg(), int64(4), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(2, 2))

	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO way_nodes (way_id, node_id, sequence_id) VALUES ($1, $2, $3), ($4, $5, $6), ($7, $8, $9), ($10, $11, $12) ON CONFLICT DO NOTHING")).
		WillReturnResult(sqlmock.NewResult(4, 4))

	waysToInsert := []models.Way{
		{ID: 3, Tags: map[string]string{"x": "1"}, NodeIDs: []int64{1, 2}},
		{ID: 4, Tags: map[string]string{"y": "2"}, NodeIDs: []int64{3, 4}},
	}

	err = repo.InsertBatches(waysToInsert, 2)
	require.NoError(t, err)

	mock.ExpectCommit()
	require.NoError(t, tx.Commit())
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestTxWayRepository_GetNearestWay(t *testing.T) {
	t.Run("returns nearest way", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer func() { mock.ExpectClose(); _ = db.Close() }()

		mock.ExpectBegin()
		tx, err := db.Begin()
		require.NoError(t, err)

		repo := NewTxWayRepository(tx)

		tags := map[string]string{"highway": "cycleway"}
		tagsJSON, _ := json.Marshal(tags)

		queryPattern := "(?s)WITH nearest_nodes.*candidate_ways.*way_geometries.*closest.*ARRAY_AGG"
		mock.ExpectQuery(queryPattern).
			WithArgs(53.5, -113.5).
			WillReturnRows(sqlmock.NewRows([]string{"id", "tags", "node_ids"}).AddRow(int64(7), tagsJSON, "{11,22}"))

		way, err := repo.GetNearestWay(53.5, -113.5)
		require.NoError(t, err)
		require.NotNil(t, way)
		require.Equal(t, int64(7), way.ID)
		require.Equal(t, tags, way.Tags)
		require.Equal(t, []int64{11, 22}, way.NodeIDs)

		mock.ExpectCommit()
		require.NoError(t, tx.Commit())
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("returns not found sentinel", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer func() { mock.ExpectClose(); _ = db.Close() }()

		mock.ExpectBegin()
		tx, err := db.Begin()
		require.NoError(t, err)

		repo := NewTxWayRepository(tx)

		queryPattern := "(?s)WITH nearest_nodes.*candidate_ways.*way_geometries.*closest.*ARRAY_AGG"
		mock.ExpectQuery(queryPattern).
			WithArgs(0.0, 0.0).
			WillReturnError(sql.ErrNoRows)

		way, err := repo.GetNearestWay(0.0, 0.0)
		require.Error(t, err)
		require.Nil(t, way)
		require.ErrorIs(t, err, repository.ErrWayNotFound)

		mock.ExpectCommit()
		require.NoError(t, tx.Commit())
		require.NoError(t, mock.ExpectationsWereMet())
	})
}
