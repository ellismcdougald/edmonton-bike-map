package geo

import (
	"testing"

	"github.com/ellismcdougald/edmonton-bike-map/internal/models"
	"github.com/stretchr/testify/require"
)

func TestMapWaysToFeatures_Success(t *testing.T) {
    nodes := map[int64]models.Node{
        1: {ID: 1, Latitude: 50.0, Longitude: -113.0},
        2: {ID: 2, Latitude: 50.1, Longitude: -113.1},
    }

    ways := []models.Way{{ID: 42, Tags: map[string]string{"highway": "cycleway"}, NodeIDs: []int64{1, 2}}}

    features, err := MapWaysToFeatures(ways, nodes)
    require.NoError(t, err)
    require.Len(t, features, 1)

    f := features[0]
    require.Equal(t, "Feature", f.Type)
    // Geometry
    require.Equal(t, "LineString", f.Geometry.Type)
    require.Len(t, f.Geometry.Coordinates, 2)
    // Coordinates are [longitude, latitude]
    require.InDelta(t, nodes[1].Longitude, f.Geometry.Coordinates[0][0], 1e-9)
    require.InDelta(t, nodes[1].Latitude, f.Geometry.Coordinates[0][1], 1e-9)

    // Properties include tags and id
    require.Equal(t, "42", f.Properties["id"])
    require.Equal(t, "cycleway", f.Properties["highway"])
}

func TestMapWaysToFeatures_MissingNode_Error(t *testing.T) {
    nodes := map[int64]models.Node{
        1: {ID: 1, Latitude: 50.0, Longitude: -113.0},
        // note: node 2 missing
    }

    ways := []models.Way{{ID: 42, Tags: map[string]string{"highway": "cycleway"}, NodeIDs: []int64{1, 2}}}

    features, err := MapWaysToFeatures(ways, nodes)
    require.Error(t, err)
    require.Nil(t, features)
}

func TestMapWaysToFeatures_EmptyInput_ReturnsEmpty(t *testing.T) {
    nodes := map[int64]models.Node{}
    ways := []models.Way{}

    features, err := MapWaysToFeatures(ways, nodes)
    require.NoError(t, err)
    // function returns nil slice when empty; accept nil or empty
    require.Empty(t, features)
}

func TestMapWaysToFeatures_SingleNodeWay(t *testing.T) {
    nodes := map[int64]models.Node{
        1: {ID: 1, Latitude: 49.0, Longitude: -113.5},
    }
    ways := []models.Way{{ID: 7, Tags: map[string]string{"name":"single"}, NodeIDs: []int64{1}}}

    features, err := MapWaysToFeatures(ways, nodes)
    require.NoError(t, err)
    require.Len(t, features, 1)
    require.Len(t, features[0].Geometry.Coordinates, 1)
    require.Equal(t, "single", features[0].Properties["name"])
}

func TestMapWaysToFeatures_ZeroNodeWay(t *testing.T) {
    nodes := map[int64]models.Node{}
    ways := []models.Way{{ID: 8, Tags: map[string]string{"name":"empty"}, NodeIDs: []int64{}}}

    features, err := MapWaysToFeatures(ways, nodes)
    require.NoError(t, err)
    require.Len(t, features, 1)
    // zero coordinates but feature still created
    require.Len(t, features[0].Geometry.Coordinates, 0)
    require.Equal(t, "empty", features[0].Properties["name"])
}

func TestMapWaysToFeatures_DuplicateNodes(t *testing.T) {
    nodes := map[int64]models.Node{
        1: {ID: 1, Latitude: 48.0, Longitude: -113.0},
    }
    ways := []models.Way{{ID: 9, Tags: map[string]string{"dup":"yes"}, NodeIDs: []int64{1, 1}}}

    features, err := MapWaysToFeatures(ways, nodes)
    require.NoError(t, err)
    require.Len(t, features, 1)
    require.Len(t, features[0].Geometry.Coordinates, 2)
    // both coordinates should equal node 1
    require.Equal(t, features[0].Geometry.Coordinates[0], features[0].Geometry.Coordinates[1])
}
