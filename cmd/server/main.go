package main

import (
	//"fmt"
	"fmt"
	"path/filepath"

	"github.com/ellismcdougald/edmonton-bike-map/pkg/data"
)

func main() {
	//fmt.Println("Main")
	fileName := filepath.Join("osm_bike_data.json")
	network, _ := data.BuildGraph(fileName)
	fmt.Println(network)
	//fmt.Println(route_data)
}