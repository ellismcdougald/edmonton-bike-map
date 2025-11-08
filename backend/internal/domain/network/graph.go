package network

import (
	"log"

	"github.com/ellismcdougald/edmonton-bike-map/internal/models"
)

// buildGraph constructs a Network graph from nodes, ways, and reviews.
func buildGraph(nodes map[int64]models.Node, ways []models.Way, reviews map[int64][]models.Review) (*models.Network, error) {
	edgesByNode := make(map[int64][]models.Edge, len(nodes))

	for _, way := range ways {
		tagsMultiplier := computeTagsMultiplier(way.Tags)
		reviewMultiplier := 1.0
		if wayReviews, ok := reviews[way.ID]; ok {
			reviewMultiplier = computeReviewMultiplier(wayReviews)
		}

		for i := 0; i < len(way.NodeIDs)-1; i++ {
			fromID := way.NodeIDs[i]
			toID := way.NodeIDs[i+1]

			fromNode, fromExists := nodes[fromID]
			toNode, toExists := nodes[toID]
			if !fromExists || !toExists {
				log.Printf("Warning: node missing for way %d: %d -> %d", way.ID, fromID, toID)
				continue
			}

			dist := haversineDistance(fromNode.Latitude, fromNode.Longitude, toNode.Latitude, toNode.Longitude)
			weight := dist * tagsMultiplier * reviewMultiplier

			edgesByNode[fromID] = append(edgesByNode[fromID], models.Edge{
				To:     toID,
				Weight: weight,
			})

			// Add reverse edge if not a oneway
			if way.Tags[tagOneway] != "yes" {
				edgesByNode[toID] = append(edgesByNode[toID], models.Edge{
					To:     fromID,
					Weight: weight,
				})
			}
		}
	}

	return &models.Network{
		Nodes: nodes,
		Edges: edgesByNode,
	}, nil
}
