package txrepo

import (
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/ellismcdougald/edmonton-bike-map/internal/models"
	"github.com/stretchr/testify/require"
)

func TestTxUserRepository_GetByUsername(t *testing.T) {
    db, mock, err := sqlmock.New()
    require.NoError(t, err)
    defer func() { mock.ExpectClose(); _ = db.Close() }()

    mock.ExpectBegin()
    tx, err := db.Begin()
    require.NoError(t, err)

    rows := sqlmock.NewRows([]string{"id", "username", "password"}).AddRow(int64(1), "alice", "pwhash")
    mock.ExpectQuery(regexp.QuoteMeta("SELECT id, username, password FROM users WHERE username = $1")).
        WithArgs("alice").
        WillReturnRows(rows)

    repo := NewTxUserRepository(tx)
    user, err := repo.GetByUsername("alice")
    require.NoError(t, err)
    require.NotNil(t, user)
    require.Equal(t, int64(1), user.ID)
    require.Equal(t, "alice", user.Username)

    mock.ExpectCommit()
    require.NoError(t, tx.Commit())
    require.NoError(t, mock.ExpectationsWereMet())
}

func TestTxUserRepository_Create(t *testing.T) {
    db, mock, err := sqlmock.New()
    require.NoError(t, err)
    defer func() { mock.ExpectClose(); _ = db.Close() }()

    mock.ExpectBegin()
    tx, err := db.Begin()
    require.NoError(t, err)

    u := &models.User{Username: "bob", Password: "secret"}
    mock.ExpectExec(regexp.QuoteMeta("INSERT INTO users (username, password) VALUES ($1, $2)")).
        WithArgs("bob", "secret").
        WillReturnResult(sqlmock.NewResult(1, 1))

    repo := NewTxUserRepository(tx)
    err = repo.Create(u)
    require.NoError(t, err)

    mock.ExpectCommit()
    require.NoError(t, tx.Commit())
    require.NoError(t, mock.ExpectationsWereMet())
}

func TestTxUserRepository_UsernameExists(t *testing.T) {
    db, mock, err := sqlmock.New()
    require.NoError(t, err)
    defer func() { mock.ExpectClose(); _ = db.Close() }()

    mock.ExpectBegin()
    tx, err := db.Begin()
    require.NoError(t, err)

    rowsTrue := sqlmock.NewRows([]string{"exists"}).AddRow(true)
    mock.ExpectQuery(regexp.QuoteMeta("SELECT EXISTS (SELECT 1 FROM users WHERE username = $1)")).
        WithArgs("carol").
        WillReturnRows(rowsTrue)

    repo := NewTxUserRepository(tx)
    ok, err := repo.UsernameExists("carol")
    require.NoError(t, err)
    require.True(t, ok)

    rowsFalse := sqlmock.NewRows([]string{"exists"}).AddRow(false)
    mock.ExpectQuery(regexp.QuoteMeta("SELECT EXISTS (SELECT 1 FROM users WHERE username = $1)")).
        WithArgs("dave").
        WillReturnRows(rowsFalse)

    ok2, err := repo.UsernameExists("dave")
    require.NoError(t, err)
    require.False(t, ok2)

    mock.ExpectCommit()
    require.NoError(t, tx.Commit())
    require.NoError(t, mock.ExpectationsWereMet())
}
