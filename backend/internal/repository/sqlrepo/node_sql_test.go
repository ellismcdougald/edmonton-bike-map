package sqlrepo

import (
	"database/sql"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/ellismcdougald/edmonton-bike-map/internal/models"
	"github.com/stretchr/testify/require"
)

func TestSQLNodeRepository_Insert(t *testing.T) {
	tests := []struct {
		name      string
		node      models.Node
		mockSetup func(mock sqlmock.Sqlmock)
		wantErr   bool
	}{
		{
			name: "successful insert",
			node: models.Node{ID: 1, Latitude: 10.0, Longitude: 20.0},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta("INSERT INTO nodes (id, latitude, longitude) VALUES ($1, $2, $3) ON CONFLICT (id) DO NOTHING")).
					WithArgs(int64(1), 10.0, 20.0).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			wantErr: false,
		},
		{
			name: "insert error",
			node: models.Node{ID: 2, Latitude: 0.0, Longitude: 0.0},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta("INSERT INTO nodes (id, latitude, longitude) VALUES ($1, $2, $3) ON CONFLICT (id) DO NOTHING")).
					WithArgs(int64(2), 0.0, 0.0).
					WillReturnError(errors.New("insert failed"))
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
				if err := db.Close(); err != nil {
					t.Errorf("failed to close db: %v", err)
				}
			}()

			tt.mockSetup(mock)

			repo := NewSQLNodeRepository(db)
			err = repo.Insert(tt.node)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestSQLNodeRepository_GetNode(t *testing.T) {
	tests := []struct {
		name      string
		nodeID    int64
		mockSetup func(mock sqlmock.Sqlmock)
		wantNode  *models.Node
		wantErr   bool
		errMsg    string
	}{
		{
			name:   "node found",
			nodeID: 1,
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "latitude", "longitude"}).
					AddRow(int64(1), 10.0, 20.0)
				mock.ExpectQuery(regexp.QuoteMeta("SELECT id, latitude, longitude FROM nodes WHERE id = $1")).
					WithArgs(int64(1)).
					WillReturnRows(rows)
			},
			wantNode: &models.Node{ID: 1, Latitude: 10.0, Longitude: 20.0},
			wantErr:  false,
		},
		{
			name:   "node not found",
			nodeID: 2,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT id, latitude, longitude FROM nodes WHERE id = $1")).
					WithArgs(int64(2)).
					WillReturnError(sql.ErrNoRows)
			},
			wantNode: nil,
			wantErr:  true,
			errMsg:   "node not found",
		},
		{
			name:   "query error",
			nodeID: 3,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT id, latitude, longitude FROM nodes WHERE id = $1")).
					WithArgs(int64(3)).
					WillReturnError(errors.New("query failed"))
			},
			wantNode: nil,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer func() {
				mock.ExpectClose()
				if err := db.Close(); err != nil {
					t.Errorf("failed to close db: %v", err)
				}
			}()

			tt.mockSetup(mock)

			repo := NewSQLNodeRepository(db)
			gotNode, err := repo.GetNode(tt.nodeID)

			if tt.wantErr {
				require.Error(t, err)
				if tt.errMsg != "" {
					require.EqualError(t, err, tt.errMsg)
				}
				require.Nil(t, gotNode)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.wantNode, gotNode)
			}

			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestSQLNodeRepository_GetAllNodes(t *testing.T) {
	tests := []struct {
		name      string
		mockSetup func(mock sqlmock.Sqlmock)
		wantNodes map[int64]models.Node
		wantErr   bool
	}{
		{
			name: "multiple nodes",
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "latitude", "longitude"}).
					AddRow(int64(1), 10.0, 20.0).
					AddRow(int64(2), 11.0, 21.0).
					AddRow(int64(3), 12.0, 22.0)
				mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT 
			id,
			latitude,
			longitude
		FROM nodes;`)).WillReturnRows(rows)
			},
			wantNodes: map[int64]models.Node{
				1: {ID: 1, Latitude: 10.0, Longitude: 20.0},
				2: {ID: 2, Latitude: 11.0, Longitude: 21.0},
				3: {ID: 3, Latitude: 12.0, Longitude: 22.0},
			},
			wantErr: false,
		},
		{
			name: "query error",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT 
			id,
			latitude,
			longitude
		FROM nodes;`)).WillReturnError(errors.New("query failed"))
			},
			wantNodes: nil,
			wantErr:   true,
		},
		{
			name: "scan error",
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "latitude", "longitude"}).
					AddRow(int64(1), 10.0, 20.0).
					AddRow(int64(2), "bad", 21.0) // latitude wrong type triggers scan error
				mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT 
			id,
			latitude,
			longitude
		FROM nodes;`)).WillReturnRows(rows)
			},
			wantNodes: nil,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer func() {
				mock.ExpectClose()
				if err := db.Close(); err != nil {
					t.Errorf("failed to close db: %v", err)
				}
			}()

			tt.mockSetup(mock)

			repo := NewSQLNodeRepository(db)
			gotNodes, err := repo.GetAllNodes()

			if tt.wantErr {
				require.Error(t, err)
				require.Nil(t, gotNodes)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.wantNodes, gotNodes)
			}

			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestSQLNodeRepository_InsertBatches(t *testing.T) {
	tests := []struct {
		name      string
		nodes     []models.Node
		batchSize int
		mockSetup func(mock sqlmock.Sqlmock)
		wantErr   bool
	}{
		{
			name:      "normal batch insert",
			nodes:     []models.Node{{ID: 1, Latitude: 1, Longitude: 1}, {ID: 2, Latitude: 2, Longitude: 2}},
			batchSize: 2,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO nodes (id, latitude, longitude) VALUES ($1, $2, $3), ($4, $5, $6) ON CONFLICT (id) DO NOTHING`)).
					WithArgs(int64(1), 1.0, 1.0, int64(2), 2.0, 2.0).
					WillReturnResult(sqlmock.NewResult(2, 2))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name:      "insert exec error",
			nodes:     []models.Node{{ID: 1, Latitude: 1, Longitude: 1}},
			batchSize: 1,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO nodes (id, latitude, longitude) VALUES ($1, $2, $3) ON CONFLICT (id) DO NOTHING`)).
					WithArgs(int64(1), 1.0, 1.0).
					WillReturnError(errors.New("insert failed"))
				mock.ExpectRollback()
			},
			wantErr: true,
		},
		{
			name:      "empty slice does nothing",
			nodes:     []models.Node{},
			batchSize: 2,
			mockSetup: func(mock sqlmock.Sqlmock) {},
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer func() {
				mock.ExpectClose()
				if err := db.Close(); err != nil {
					t.Errorf("failed to close db: %v", err)
				}
			}()

			tt.mockSetup(mock)

			repo := NewSQLNodeRepository(db)
			err = repo.InsertBatches(tt.nodes, tt.batchSize)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
