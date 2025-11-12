package sqlrepo

import (
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/ellismcdougald/edmonton-bike-map/internal/models"
	"github.com/stretchr/testify/require"
)

func TestSQLReviewRepository_CreateReview(t *testing.T) {
    db, mock, err := sqlmock.New()
    require.NoError(t, err)
    defer func() {
        mock.ExpectClose()
        _ = db.Close()
    }()

    rev := &models.Review{WayID: 1, UserID: 2, Rating: 5, Comment: "nice"}

    mock.ExpectExec(regexp.QuoteMeta("INSERT INTO reviews (\n            way_id,\n            user_id,\n            rating,\n            comment\n        ) VALUES ($1, $2, $3, $4)")).
        WithArgs(int64(1), int64(2), 5, "nice").
        WillReturnResult(sqlmock.NewResult(1, 1))

    repo := NewSQLReviewRepository(db)
    err = repo.CreateReview(rev)
    require.NoError(t, err)
    require.NoError(t, mock.ExpectationsWereMet())
}

func TestSQLReviewRepository_GetReviews(t *testing.T) {
    db, mock, err := sqlmock.New()
    require.NoError(t, err)
    defer func() {
        mock.ExpectClose()
        _ = db.Close()
    }()

    created := time.Now()
    rows := sqlmock.NewRows([]string{"way_id", "user_id", "rating", "comment", "created_at", "username"}).
        AddRow(int64(1), int64(2), 4, "ok", created, "bob")

    mock.ExpectQuery(regexp.QuoteMeta("SELECT\n            r.way_id,\n            r.user_id,\n            r.rating,\n            r.comment,\n            r.created_at,\n            u.username\n        FROM reviews r\n        LEFT JOIN users u ON r.user_id = u.id\n        WHERE r.way_id = $1;" )).
        WithArgs(int64(1)).
        WillReturnRows(rows)

    repo := NewSQLReviewRepository(db)
    reviews, err := repo.GetReviews(1)
    require.NoError(t, err)
    require.Len(t, reviews, 1)
    require.Equal(t, int64(1), reviews[0].WayID)
    require.Equal(t, int64(2), reviews[0].UserID)
    require.Equal(t, "bob", reviews[0].Username)

    require.NoError(t, mock.ExpectationsWereMet())
}

func TestSQLReviewRepository_GetAllReviews(t *testing.T) {
    db, mock, err := sqlmock.New()
    require.NoError(t, err)
    defer func() {
        mock.ExpectClose()
        _ = db.Close()
    }()

    created := time.Now()
    rows := sqlmock.NewRows([]string{"way_id", "user_id", "rating", "comment", "created_at", "username"}).
        AddRow(int64(1), int64(2), 4, "ok", created, "bob").
        AddRow(int64(2), int64(3), 5, "great", created, "alice")

    mock.ExpectQuery(regexp.QuoteMeta("SELECT\n            r.way_id,\n            r.user_id,\n            r.rating,\n            r.comment,\n            r.created_at,\n            u.username\n        FROM reviews r\n        LEFT JOIN users u ON r.user_id = u.id;" )).
        WillReturnRows(rows)

    repo := NewSQLReviewRepository(db)
    res, err := repo.GetAllReviews()
    require.NoError(t, err)
    require.Len(t, res, 2)
    require.Equal(t, int64(1), res[int64(1)][0].WayID)
    require.Equal(t, int64(2), res[int64(2)][0].WayID)

    require.NoError(t, mock.ExpectationsWereMet())
}

func TestSQLReviewRepository_InsertBatches(t *testing.T) {
    db, mock, err := sqlmock.New()
    require.NoError(t, err)
    defer func() { mock.ExpectClose(); _ = db.Close() }()

    repo := NewSQLReviewRepository(db)

    // normal batch of two reviews
    mock.ExpectExec(regexp.QuoteMeta("INSERT INTO reviews (way_id, user_id, rating, comment) VALUES ($1, $2, $3, $4), ($5, $6, $7, $8)")).
        WithArgs(int64(10), int64(2), 5, "great", int64(11), int64(3), 4, "ok").
        WillReturnResult(sqlmock.NewResult(2, 2))

    reviews := []models.Review{{WayID: 10, UserID: 2, Rating: 5, Comment: "great"}, {WayID: 11, UserID: 3, Rating: 4, Comment: "ok"}}
    reqErr := repo.InsertBatches(reviews, 2)
    require.NoError(t, reqErr)

    // empty slice does nothing
    reqErr = repo.InsertBatches([]models.Review{}, 2)
    require.NoError(t, reqErr)

    // exec error
    mock.ExpectExec(regexp.QuoteMeta("INSERT INTO reviews (way_id, user_id, rating, comment) VALUES ($1, $2, $3, $4)")).
        WithArgs(int64(12), int64(4), 3, "meh").
        WillReturnError(fmt.Errorf("insert failure"))

    reqErr = repo.InsertBatches([]models.Review{{WayID: 12, UserID: 4, Rating: 3, Comment: "meh"}}, 1)
    require.Error(t, reqErr)

    require.NoError(t, mock.ExpectationsWereMet())
}
