package network

import (
	"github.com/ellismcdougald/edmonton-bike-map/internal/models"
	"github.com/ellismcdougald/edmonton-bike-map/internal/service"
)

func BuildNetwork(nodeService service.NodeService, wayService service.WayService) (*models.Network, error) {
	allNodes, err := nodeService.GetAllNodes()
	if err != nil {
		return nil, err
	}

	allWays, err := wayService.GetAllWays()
	if err != nil {
		return nil, err
	}
	return buildGraph(allNodes, allWays)
}
