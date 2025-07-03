package main

import (
	"log"
	"net/http"
	"path/filepath"

	"github.com/ellismcdougald/edmonton-bike-map/pkg/data"
	"github.com/ellismcdougald/edmonton-bike-map/pkg/server"
)

func main() {
	fileName := filepath.Join("osm_bike_data.json")
	network, _ := data.BuildGraph(fileName)

	mux := http.NewServeMux()
	server.RegisterRoutes(mux, network)

	addr := ":8080"
	log.Printf("Starting server on %s\n", addr)
	err := http.ListenAndServe(addr, mux)
	if err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}