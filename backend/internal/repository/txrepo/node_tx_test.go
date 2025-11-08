package txrepo

import (
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/ellismcdougald/edmonton-bike-map/internal/models"
	"github.com/stretchr/testify/require"
)

func TestTxNodeRepository_Insert(t *testing.T) {
    db, mock, err := sqlmock.New()
    require.NoError(t, err)
    defer func() { mock.ExpectClose(); _ = db.Close() }()

    mock.ExpectBegin()
    tx, err := db.Begin()
    require.NoError(t, err)

    repo := NewTxNodeRepository(tx)

    mock.ExpectExec(regexp.QuoteMeta("INSERT INTO nodes (id, latitude, longitude) VALUES ($1, $2, $3) ON CONFLICT (id) DO NOTHING")).
        WithArgs(int64(1), 10.0, 20.0).
        WillReturnResult(sqlmock.NewResult(1, 1))

    err = repo.Insert(models.Node{ID: 1, Latitude: 10.0, Longitude: 20.0})
    require.NoError(t, err)

    mock.ExpectCommit()
    require.NoError(t, tx.Commit())
    require.NoError(t, mock.ExpectationsWereMet())
}

func TestTxNodeRepository_GetNode(t *testing.T) {
    db, mock, err := sqlmock.New()
    require.NoError(t, err)
    defer func() { mock.ExpectClose(); _ = db.Close() }()

    mock.ExpectBegin()
    tx, err := db.Begin()
    require.NoError(t, err)

    repo := NewTxNodeRepository(tx)

    rows := sqlmock.NewRows([]string{"id", "latitude", "longitude"}).AddRow(int64(1), 10.0, 20.0)
    mock.ExpectQuery(regexp.QuoteMeta("SELECT id, latitude, longitude FROM nodes WHERE id = $1")).WithArgs(int64(1)).WillReturnRows(rows)

    n, err := repo.GetNode(1)
    require.NoError(t, err)
    require.Equal(t, int64(1), n.ID)

    mock.ExpectCommit()
    require.NoError(t, tx.Commit())
    require.NoError(t, mock.ExpectationsWereMet())
}

func TestTxNodeRepository_GetAllNodes(t *testing.T) {
    db, mock, err := sqlmock.New()
    require.NoError(t, err)
    defer func() { mock.ExpectClose(); _ = db.Close() }()

    mock.ExpectBegin()
    tx, err := db.Begin()
    require.NoError(t, err)

    repo := NewTxNodeRepository(tx)

    allRows := sqlmock.NewRows([]string{"id", "latitude", "longitude"}).
        AddRow(int64(1), 10.0, 20.0).
        AddRow(int64(2), 11.0, 21.0)
    mock.ExpectQuery(regexp.QuoteMeta("SELECT id, latitude, longitude FROM nodes")).WillReturnRows(allRows)

    nodes, err := repo.GetAllNodes()
    require.NoError(t, err)
    require.Len(t, nodes, 2)

    mock.ExpectCommit()
    require.NoError(t, tx.Commit())
    require.NoError(t, mock.ExpectationsWereMet())
}

func TestTxNodeRepository_InsertBatches(t *testing.T) {
    db, mock, err := sqlmock.New()
    require.NoError(t, err)
    defer func() { mock.ExpectClose(); _ = db.Close() }()

    mock.ExpectBegin()
    tx, err := db.Begin()
    require.NoError(t, err)

    repo := NewTxNodeRepository(tx)

    mock.ExpectExec(regexp.QuoteMeta("INSERT INTO nodes (id, latitude, longitude) VALUES ($1, $2, $3), ($4, $5, $6) ON CONFLICT (id) DO NOTHING")).
        WithArgs(int64(3), 1.0, 1.0, int64(4), 2.0, 2.0).
        WillReturnResult(sqlmock.NewResult(2, 2))

    err = repo.InsertBatches([]models.Node{{ID: 3, Latitude: 1, Longitude: 1}, {ID: 4, Latitude: 2, Longitude: 2}}, 2)
    require.NoError(t, err)

    mock.ExpectCommit()
    require.NoError(t, tx.Commit())
    require.NoError(t, mock.ExpectationsWereMet())
}
