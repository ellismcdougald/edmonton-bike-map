package network

import (
	"math"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHaversineDistance_ZeroAndSymmetric(t *testing.T) {
	d1 := haversineDistance(50.0, -113.0, 50.0, -113.0)
	require.Equal(t, 0.0, d1)

	d2 := haversineDistance(50.0, -113.0, 51.0, -114.0)
	d3 := haversineDistance(51.0, -114.0, 50.0, -113.0)
	require.InDelta(t, d2, d3, 1e-6)
}

func TestHaversineDistance_OneDegreeAtEquator(t *testing.T) {
	// distance for 1 degree at equator: R * (pi/180)
	expected := 6371000.0 * (math.Pi / 180.0)
	got := haversineDistance(0.0, 0.0, 0.0, 1.0)
	// allow a small absolute tolerance (1 millimeter)
	require.InDelta(t, expected, got, 1e-3)
}

func TestHaversineDistance_SmallDeltaLatitude(t *testing.T) {
	// 0.0001 degree latitude ~ 11.1195 meters
	got := haversineDistance(50.0, -113.0, 50.0001, -113.0)
	require.Greater(t, got, 11.0)
	require.Less(t, got, 12.0)
}
