package model

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
)

func TestDBNodeStore_Insert(t *testing.T) {
	tests := []struct {
		name      string
		node      DBNode
		mockSetup func(mock sqlmock.Sqlmock)
		wantErr   bool
	}{
		{
			name: "successful insert",
			node: DBNode{ID: 1, Latitude: 10.0, Longitude: 20.0},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`INSERT INTO nodes \(id, latitude, longitude\) VALUES \(\$1, \$2, \$3\) ON CONFLICT \(id\) DO NOTHING`).
					WithArgs(int64(1), 10.0, 20.0).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			wantErr: false,
		},
		{
			name: "insert error",
			node: DBNode{ID: 2, Latitude: 0.0, Longitude: 0.0},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`INSERT INTO nodes \(id, latitude, longitude\) VALUES \(\$1, \$2, \$3\) ON CONFLICT \(id\) DO NOTHING`).
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
			defer db.Close()

			tt.mockSetup(mock)

			store := &DBNodeStore{DB: db}
			err = store.Insert(tt.node)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestDBNodeStore_GetNode(t *testing.T) {
	tests := []struct {
		name      string
		nodeID    int64
		mockSetup func(mock sqlmock.Sqlmock)
		wantNode  *DBNode
		wantErr   bool
		errMsg    string
	}{
		{
			name:   "node found",
			nodeID: 1,
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "latitude", "longitude"}).
					AddRow(int64(1), 10.0, 20.0)
				mock.ExpectQuery(`SELECT id, latitude, longitude FROM nodes WHERE id = \$1`).
					WithArgs(int64(1)).
					WillReturnRows(rows)
			},
			wantNode: &DBNode{ID: 1, Latitude: 10.0, Longitude: 20.0},
			wantErr:  false,
		},
		{
			name:   "node not found",
			nodeID: 2,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT id, latitude, longitude FROM nodes WHERE id = \$1`).
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
				mock.ExpectQuery(`SELECT id, latitude, longitude FROM nodes WHERE id = \$1`).
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
			defer db.Close()

			tt.mockSetup(mock)

			store := &DBNodeStore{DB: db}
			gotNode, err := store.GetNode(tt.nodeID)

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

func TestDBNodeStore_GetAllNodes(t *testing.T) {
	tests := []struct {
		name      string
		mockSetup func(mock sqlmock.Sqlmock)
		wantNodes map[int64]DBNode
		wantErr   bool
	}{
		{
			name: "multiple nodes",
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "latitude", "longitude"}).
					AddRow(int64(1), 10.0, 20.0).
					AddRow(int64(2), 11.0, 21.0).
					AddRow(int64(3), 12.0, 22.0)
				mock.ExpectQuery(`SELECT 
			id,
			latitude,
			longitude
		FROM nodes;`).WillReturnRows(rows)
			},
			wantNodes: map[int64]DBNode{
				1: {ID: 1, Latitude: 10.0, Longitude: 20.0},
				2: {ID: 2, Latitude: 11.0, Longitude: 21.0},
				3: {ID: 3, Latitude: 12.0, Longitude: 22.0},
			},
			wantErr: false,
		},
		{
			name: "query error",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT 
			id,
			latitude,
			longitude
		FROM nodes;`).WillReturnError(errors.New("query failed"))
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
				mock.ExpectQuery(`SELECT 
			id,
			latitude,
			longitude
		FROM nodes;`).WillReturnRows(rows)
			},
			wantNodes: nil,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			tt.mockSetup(mock)

			store := &DBNodeStore{DB: db}
			gotNodes, err := store.GetAllNodes()

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
