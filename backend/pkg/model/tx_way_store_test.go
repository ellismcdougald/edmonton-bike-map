package model

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
)

func TestTxWayStore_Insert(t *testing.T) {
	tests := []struct {
		name      string
		way       DBWay
		mockSetup func(mock sqlmock.Sqlmock)
		wantErr   bool
	}{
		{
			name: "successful insert with nodes",
			way: DBWay{
				ID:      1,
				Tags:    map[string]string{"highway": "primary", "name": "Main St"},
				NodeIDs: []int64{10, 20, 30},
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				// Insert way
				mock.ExpectExec(`INSERT INTO ways \(id, tags\) VALUES \(\$1, \$2\) ON CONFLICT \(id\) DO NOTHING`).
					WithArgs(int64(1), []byte(`{"highway":"primary","name":"Main St"}`)).
					WillReturnResult(sqlmock.NewResult(1, 1))
				// Insert way_nodes
				mock.ExpectExec(`INSERT INTO way_nodes \(way_id, node_id, sequence_id\) VALUES \(\$1, \$2, \$3\) ON CONFLICT DO NOTHING`).
					WithArgs(int64(1), int64(10), 0).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec(`INSERT INTO way_nodes \(way_id, node_id, sequence_id\) VALUES \(\$1, \$2, \$3\) ON CONFLICT DO NOTHING`).
					WithArgs(int64(1), int64(20), 1).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec(`INSERT INTO way_nodes \(way_id, node_id, sequence_id\) VALUES \(\$1, \$2, \$3\) ON CONFLICT DO NOTHING`).
					WithArgs(int64(1), int64(30), 2).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			wantErr: false,
		},
		{
			name: "successful insert with empty tags and no nodes",
			way: DBWay{
				ID:      2,
				Tags:    map[string]string{},
				NodeIDs: []int64{},
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(`INSERT INTO ways \(id, tags\) VALUES \(\$1, \$2\) ON CONFLICT \(id\) DO NOTHING`).
					WithArgs(int64(2), []byte(`{}`)).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			wantErr: false,
		},
		{
			name: "way insert error",
			way: DBWay{
				ID:      3,
				Tags:    map[string]string{"highway": "secondary"},
				NodeIDs: []int64{40},
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(`INSERT INTO ways \(id, tags\) VALUES \(\$1, \$2\) ON CONFLICT \(id\) DO NOTHING`).
					WithArgs(int64(3), []byte(`{"highway":"secondary"}`)).
					WillReturnError(errors.New("way insert failed"))
			},
			wantErr: true,
		},
		{
			name: "way_node insert error",
			way: DBWay{
				ID:      4,
				Tags:    map[string]string{"highway": "tertiary"},
				NodeIDs: []int64{50},
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(`INSERT INTO ways \(id, tags\) VALUES \(\$1, \$2\) ON CONFLICT \(id\) DO NOTHING`).
					WithArgs(int64(4), []byte(`{"highway":"tertiary"}`)).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec(`INSERT INTO way_nodes \(way_id, node_id, sequence_id\) VALUES \(\$1, \$2, \$3\) ON CONFLICT DO NOTHING`).
					WithArgs(int64(4), int64(50), 0).
					WillReturnError(errors.New("way_node insert failed"))
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

			tx, err := db.Begin()
			require.NoError(t, err)

			store := &TxWayStore{Tx: tx}
			err = store.Insert(tt.way)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestTxWayStore_InsertBatch(t *testing.T) {
	tests := []struct {
		name      string
		ways      []DBWay
		mockSetup func(mock sqlmock.Sqlmock)
		wantErr   bool
	}{
		{
			name: "successful batch insert",
			ways: []DBWay{
				{ID: 1, Tags: map[string]string{"highway": "primary"}, NodeIDs: []int64{10, 20}},
				{ID: 2, Tags: map[string]string{"highway": "secondary"}, NodeIDs: []int64{30}},
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				// Insert ways
				mock.ExpectExec(`INSERT INTO ways \(id, tags\) VALUES \(\$1, \$2\), \(\$3, \$4\) ON CONFLICT \(id\) DO NOTHING`).
					WithArgs(int64(1), []byte(`{"highway":"primary"}`), int64(2), []byte(`{"highway":"secondary"}`)).
					WillReturnResult(sqlmock.NewResult(2, 2))
				// Insert way_nodes
				mock.ExpectExec(`INSERT INTO way_nodes \(way_id, node_id, sequence_id\) VALUES \(\$1, \$2, \$3\), \(\$4, \$5, \$6\), \(\$7, \$8, \$9\) ON CONFLICT DO NOTHING`).
					WithArgs(int64(1), int64(10), 0, int64(1), int64(20), 1, int64(2), int64(30), 0).
					WillReturnResult(sqlmock.NewResult(3, 3))
			},
			wantErr: false,
		},
		{
			name: "empty slice does nothing",
			ways: []DBWay{},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
			},
			wantErr: false,
		},
		{
			name: "way insert error",
			ways: []DBWay{
				{ID: 1, Tags: map[string]string{"highway": "primary"}, NodeIDs: []int64{10}},
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(`INSERT INTO ways \(id, tags\) VALUES \(\$1, \$2\) ON CONFLICT \(id\) DO NOTHING`).
					WithArgs(int64(1), []byte(`{"highway":"primary"}`)).
					WillReturnError(errors.New("way insert failed"))
			},
			wantErr: true,
		},
		{
			name: "way_nodes insert error",
			ways: []DBWay{
				{ID: 1, Tags: map[string]string{"highway": "primary"}, NodeIDs: []int64{10}},
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(`INSERT INTO ways \(id, tags\) VALUES \(\$1, \$2\) ON CONFLICT \(id\) DO NOTHING`).
					WithArgs(int64(1), []byte(`{"highway":"primary"}`)).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec(`INSERT INTO way_nodes \(way_id, node_id, sequence_id\) VALUES \(\$1, \$2, \$3\) ON CONFLICT DO NOTHING`).
					WithArgs(int64(1), int64(10), 0).
					WillReturnError(errors.New("way_nodes insert failed"))
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

			tx, err := db.Begin()
			require.NoError(t, err)

			store := &TxWayStore{Tx: tx}
			err = store.InsertBatch(tt.ways)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestTxWayStore_InsertBatchChunks(t *testing.T) {
	tests := []struct {
		name      string
		ways      []DBWay
		batchSize int
		mockSetup func(mock sqlmock.Sqlmock)
		wantErr   bool
	}{
		{
			name: "multiple chunks",
			ways: []DBWay{
				{ID: 1, Tags: map[string]string{"highway": "primary"}, NodeIDs: []int64{10}},
				{ID: 2, Tags: map[string]string{"highway": "secondary"}, NodeIDs: []int64{20}},
				{ID: 3, Tags: map[string]string{"highway": "tertiary"}, NodeIDs: []int64{30}},
			},
			batchSize: 2,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				// First batch (2 ways)
				mock.ExpectExec(`INSERT INTO ways \(id, tags\) VALUES \(\$1, \$2\), \(\$3, \$4\) ON CONFLICT \(id\) DO NOTHING`).
					WithArgs(int64(1), []byte(`{"highway":"primary"}`), int64(2), []byte(`{"highway":"secondary"}`)).
					WillReturnResult(sqlmock.NewResult(2, 2))
				mock.ExpectExec(`INSERT INTO way_nodes \(way_id, node_id, sequence_id\) VALUES \(\$1, \$2, \$3\), \(\$4, \$5, \$6\) ON CONFLICT DO NOTHING`).
					WithArgs(int64(1), int64(10), 0, int64(2), int64(20), 0).
					WillReturnResult(sqlmock.NewResult(2, 2))
				// Second batch (1 way)
				mock.ExpectExec(`INSERT INTO ways \(id, tags\) VALUES \(\$1, \$2\) ON CONFLICT \(id\) DO NOTHING`).
					WithArgs(int64(3), []byte(`{"highway":"tertiary"}`)).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec(`INSERT INTO way_nodes \(way_id, node_id, sequence_id\) VALUES \(\$1, \$2, \$3\) ON CONFLICT DO NOTHING`).
					WithArgs(int64(3), int64(30), 0).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			wantErr: false,
		},
		{
			name:      "default batch size when <= 0",
			ways:      []DBWay{{ID: 1, Tags: map[string]string{"highway": "primary"}, NodeIDs: []int64{10}}},
			batchSize: 0, // Should default to 500
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(`INSERT INTO ways \(id, tags\) VALUES \(\$1, \$2\) ON CONFLICT \(id\) DO NOTHING`).
					WithArgs(int64(1), []byte(`{"highway":"primary"}`)).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec(`INSERT INTO way_nodes \(way_id, node_id, sequence_id\) VALUES \(\$1, \$2, \$3\) ON CONFLICT DO NOTHING`).
					WithArgs(int64(1), int64(10), 0).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			wantErr: false,
		},
		{
			name:      "empty slice does nothing",
			ways:      []DBWay{},
			batchSize: 2,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
			},
			wantErr: false,
		},
		{
			name:      "insert batch error propagates",
			ways:      []DBWay{{ID: 1, Tags: map[string]string{"highway": "primary"}, NodeIDs: []int64{10}}},
			batchSize: 2,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(`INSERT INTO ways \(id, tags\) VALUES \(\$1, \$2\) ON CONFLICT \(id\) DO NOTHING`).
					WithArgs(int64(1), []byte(`{"highway":"primary"}`)).
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

			tx, err := db.Begin()
			require.NoError(t, err)

			store := &TxWayStore{Tx: tx}
			err = store.InsertBatchChunks(tt.ways, tt.batchSize)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestTxWayStore_InsertBatchDynamic(t *testing.T) {
	tests := []struct {
		name      string
		ways      []DBWay
		mockSetup func(mock sqlmock.Sqlmock)
		wantErr   bool
	}{
		{
			name: "single batch under limit",
			ways: []DBWay{
				{ID: 1, Tags: map[string]string{"highway": "primary"}, NodeIDs: []int64{10, 20}},
				{ID: 2, Tags: map[string]string{"highway": "secondary"}, NodeIDs: []int64{30}},
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				// All ways fit in one batch
				mock.ExpectExec(`INSERT INTO ways \(id, tags\) VALUES \(\$1, \$2\), \(\$3, \$4\) ON CONFLICT \(id\) DO NOTHING`).
					WithArgs(int64(1), []byte(`{"highway":"primary"}`), int64(2), []byte(`{"highway":"secondary"}`)).
					WillReturnResult(sqlmock.NewResult(2, 2))
				mock.ExpectExec(`INSERT INTO way_nodes \(way_id, node_id, sequence_id\) VALUES \(\$1, \$2, \$3\), \(\$4, \$5, \$6\), \(\$7, \$8, \$9\) ON CONFLICT DO NOTHING`).
					WithArgs(int64(1), int64(10), 0, int64(1), int64(20), 1, int64(2), int64(30), 0).
					WillReturnResult(sqlmock.NewResult(3, 3))
			},
			wantErr: false,
		},
		{
			name: "empty slice does nothing",
			ways: []DBWay{},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
			},
			wantErr: false,
		},
		{
			name: "flush error propagates",
			ways: []DBWay{
				{ID: 1, Tags: map[string]string{"highway": "primary"}, NodeIDs: []int64{10}},
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(`INSERT INTO ways \(id, tags\) VALUES \(\$1, \$2\) ON CONFLICT \(id\) DO NOTHING`).
					WithArgs(int64(1), []byte(`{"highway":"primary"}`)).
					WillReturnError(errors.New("flush failed"))
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

			tx, err := db.Begin()
			require.NoError(t, err)

			store := &TxWayStore{Tx: tx}
			err = store.InsertBatchDynamic(tt.ways)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}