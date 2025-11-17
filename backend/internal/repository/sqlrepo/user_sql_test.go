package sqlrepo

import (
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/ellismcdougald/edmonton-bike-map/internal/models"
	"github.com/stretchr/testify/require"
)

func TestSQLUserRepository_GetByUsername(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer func() {
		mock.ExpectClose()
		_ = db.Close()
	}()

	rows := sqlmock.NewRows([]string{"id", "username", "password", "cycling_speed"}).AddRow(int64(1), "alice", "pwhash", 15)

	// match the SELECT that retrieves user by username
	mock.ExpectQuery(regexp.QuoteMeta("SELECT")).WithArgs("alice").WillReturnRows(rows)

	repo := NewSQLUserRepository(db)
	user, err := repo.GetByUsername("alice")
	require.NoError(t, err)
	require.NotNil(t, user)
	require.Equal(t, int64(1), user.ID)
	require.Equal(t, "alice", user.Username)

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestSQLUserRepository_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer func() {
		mock.ExpectClose()
		_ = db.Close()
	}()

	u := &models.User{Username: "bob", Password: "secret"}

	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO users (username, password) VALUES ($1, $2)")).
		WithArgs("bob", "secret").
		WillReturnResult(sqlmock.NewResult(1, 1))

	repo := NewSQLUserRepository(db)
	err = repo.Create(u)
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestSQLUserRepository_GetByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer func() {
		mock.ExpectClose()
		_ = db.Close()
	}()

	rows := sqlmock.NewRows([]string{"id", "username", "password", "cycling_speed"}).AddRow(int64(42), "alice", "pwhash", 12)

	// match the SELECT that retrieves user by ID
	mock.ExpectQuery(regexp.QuoteMeta("SELECT")).WithArgs(int64(42)).WillReturnRows(rows)

	repo := NewSQLUserRepository(db)
	user, err := repo.GetByID(42)
	require.NoError(t, err)
	require.NotNil(t, user)
	require.Equal(t, int64(42), user.ID)
	require.Equal(t, "alice", user.Username)

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestSQLUserRepository_UsernameExists(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer func() {
		mock.ExpectClose()
		_ = db.Close()
	}()

	// true case
	rowsTrue := sqlmock.NewRows([]string{"exists"}).AddRow(true)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT EXISTS")).
		WithArgs("carol").
		WillReturnRows(rowsTrue)

	repo := NewSQLUserRepository(db)
	ok, err := repo.UsernameExists("carol")
	require.NoError(t, err)
	require.True(t, ok)

	// set expectation for false case
	rowsFalse := sqlmock.NewRows([]string{"exists"}).AddRow(false)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT EXISTS")).
		WithArgs("dave").
		WillReturnRows(rowsFalse)

	ok2, err := repo.UsernameExists("dave")
	require.NoError(t, err)
	require.False(t, ok2)

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestSQLUserRepository_UpdatePassword(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer func() {
		mock.ExpectClose()
		_ = db.Close()
	}()

	mock.ExpectExec(regexp.QuoteMeta("UPDATE users SET password = $1 WHERE id = $2")).
		WithArgs("newhash", int64(123)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	repo := NewSQLUserRepository(db)
	err = repo.UpdatePassword(123, "newhash")
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}
