package sqlrepo

import (
	"encoding/json"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/ellismcdougald/edmonton-bike-map/internal/models"
	"github.com/stretchr/testify/require"
)

func TestSQLWayRepository_Insert(t *testing.T) {
    tests := []struct {
        name      string
        way       models.Way
        mockSetup func(mock sqlmock.Sqlmock)
        wantErr   bool
    }{
        {
            name: "successful insert",
            way:  models.Way{ID: 1, Tags: map[string]string{"highway": "cycleway"}, NodeIDs: []int64{10, 20}},
            mockSetup: func(mock sqlmock.Sqlmock) {
                tagsJSON, _ := json.Marshal(map[string]string{"highway": "cycleway"})
                mock.ExpectBegin()
                mock.ExpectExec(regexp.QuoteMeta("INSERT INTO ways (id, tags) VALUES ($1, $2) ON CONFLICT (id) DO NOTHING")).
                    WithArgs(int64(1), tagsJSON).
                    WillReturnResult(sqlmock.NewResult(1, 1))

                // Expect inserts into way_nodes for each node
                mock.ExpectExec(regexp.QuoteMeta("INSERT INTO way_nodes (way_id, node_id, sequence_id) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING")).
                    WithArgs(int64(1), int64(10), 0).
                    WillReturnResult(sqlmock.NewResult(1, 1))
                mock.ExpectExec(regexp.QuoteMeta("INSERT INTO way_nodes (way_id, node_id, sequence_id) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING")).
                    WithArgs(int64(1), int64(20), 1).
                    WillReturnResult(sqlmock.NewResult(1, 1))

                mock.ExpectCommit()
            },
            wantErr: false,
        },
        {
            name: "insert ways exec error",
            way:  models.Way{ID: 2, Tags: map[string]string{}, NodeIDs: []int64{}},
            mockSetup: func(mock sqlmock.Sqlmock) {
                mock.ExpectBegin()
                mock.ExpectExec(regexp.QuoteMeta("INSERT INTO ways (id, tags) VALUES ($1, $2) ON CONFLICT (id) DO NOTHING")).
                    WithArgs(int64(2), sqlmock.AnyArg()).
                    WillReturnError(errors.New("insert failed"))
                mock.ExpectRollback()
            },
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            db, mock, err := sqlmock.New()
            require.NoError(t, err)
            defer func() {
                mock.ExpectClose()
                _ = db.Close()
            }()

            tt.mockSetup(mock)

            repo := NewSQLWayRepository(db)
            err = repo.Insert(tt.way)

            if tt.wantErr {
                require.Error(t, err)
            } else {
                require.NoError(t, err)
            }

            require.NoError(t, mock.ExpectationsWereMet())
        })
    }
}

func TestSQLWayRepository_GetWay(t *testing.T) {
    db, mock, err := sqlmock.New()
    require.NoError(t, err)
    defer func() {
        mock.ExpectClose()
        _ = db.Close()
    }()

    tags := map[string]string{"highway": "cycleway"}
    tagsJSON, _ := json.Marshal(tags)

    rows := sqlmock.NewRows([]string{"id", "tags"}).AddRow(int64(1), tagsJSON)
    mock.ExpectQuery(regexp.QuoteMeta("SELECT id, tags FROM ways WHERE id = $1")).WithArgs(int64(1)).WillReturnRows(rows)

    nodeRows := sqlmock.NewRows([]string{"node_id"}).AddRow(int64(10)).AddRow(int64(20))
    mock.ExpectQuery(regexp.QuoteMeta("SELECT node_id FROM way_nodes WHERE way_id = $1 ORDER BY sequence_id")).WithArgs(int64(1)).WillReturnRows(nodeRows)

    repo := NewSQLWayRepository(db)
    way, err := repo.GetWay(1)
    require.NoError(t, err)
    require.NotNil(t, way)
    require.Equal(t, int64(1), way.ID)
    require.Equal(t, tags, way.Tags)
    require.Equal(t, []int64{10, 20}, way.NodeIDs)

    require.NoError(t, mock.ExpectationsWereMet())
}

func TestSQLWayRepository_GetAllWays(t *testing.T) {
    db, mock, err := sqlmock.New()
    require.NoError(t, err)
    defer func() {
        mock.ExpectClose()
        _ = db.Close()
    }()

    query := regexp.QuoteMeta("SELECT\n\t\tw.id,\n\t\tw.tags,\n\t\tARRAY_AGG(wn.node_id ORDER BY wn.sequence_id) AS node_ids\n\tFROM ways w\n\tJOIN way_nodes wn ON wn.way_id = w.id\n\tGROUP BY w.id, w.tags;")

    tagsJSON1, _ := json.Marshal(map[string]string{"a": "1"})
    tagsJSON2, _ := json.Marshal(map[string]string{"b": "2"})

    // For the node_ids column use the PostgreSQL array text format which pq can parse (e.g. "{1,2}")
    rows := sqlmock.NewRows([]string{"id", "tags", "node_ids"}).
        AddRow(int64(1), tagsJSON1, "{1,2}").
        AddRow(int64(2), tagsJSON2, "{3,4}")

    mock.ExpectQuery(query).WillReturnRows(rows)

    repo := NewSQLWayRepository(db)
    ways, err := repo.GetAllWays()
    require.NoError(t, err)
    require.Len(t, ways, 2)

    require.NoError(t, mock.ExpectationsWereMet())
}

// (removed) pqArrayInt64 helper was unused and created a linter warning; tests use string array format instead.
