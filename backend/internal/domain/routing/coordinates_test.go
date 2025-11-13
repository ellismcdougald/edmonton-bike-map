package routing

import (
	"testing"

	"github.com/ellismcdougald/edmonton-bike-map/internal/models"
	"github.com/stretchr/testify/require"
)

func TestNearestNode(t *testing.T) {
	nodes := map[int64]models.Node{
		1: {ID: 1, Latitude: 53.5461, Longitude: -113.4938}, // Edmonton
		2: {ID: 2, Latitude: 51.0447, Longitude: -114.0719}, // Calgary
		3: {ID: 3, Latitude: 49.2827, Longitude: -123.1207}, // Vancouver
	}

	tests := []struct {
		name       string
		latitude   float64
		longitude  float64
		expectedID int64
	}{
		{
			name:       "Closest to Edmonton",
			latitude:   53.55,
			longitude:  -113.49,
			expectedID: 1,
		},
		{
			name:       "Closest to Calgary",
			latitude:   51.05,
			longitude:  -114.07,
			expectedID: 2,
		},
		{
			name:       "Closest to Vancouver",
			latitude:   49.28,
			longitude:  -123.12,
			expectedID: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := nearestNode(tt.latitude, tt.longitude, nodes)
			require.Equal(t, tt.expectedID, result)
		})
	}
}

func TestSquaredEucDistance(t *testing.T) {
	tests := []struct {
		name     string
		lat1     float64
		lon1     float64
		lat2     float64
		lon2     float64
		expected float64
	}{
		{
			name:     "Zero distance",
			lat1:     0,
			lon1:     0,
			lat2:     0,
			lon2:     0,
			expected: 0,
		},
		{
			name:     "One degree difference",
			lat1:     0,
			lon1:     0,
			lat2:     1,
			lon2:     0,
			expected: 1,
		},
		{
			name:     "Diagonal distance",
			lat1:     0,
			lon1:     0,
			lat2:     1,
			lon2:     1,
			expected: 2,
		},
		{
			name:     "Negative coordinates",
			lat1:     -1,
			lon1:     -1,
			lat2:     1,
			lon2:     1,
			expected: 8,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := squaredEucDistance(tt.lat1, tt.lon1, tt.lat2, tt.lon2)
			require.Equal(t, tt.expected, result)
		})
	}
}
