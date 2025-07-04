package main

import (
	"log"
	"net/http"
	"path/filepath"

	"github.com/ellismcdougald/edmonton-bike-map/pkg/data"
	"github.com/ellismcdougald/edmonton-bike-map/pkg/server"
)

func main() {
	/*
	var query string = `
	[out:json][timeout:30];

	// Find Canada area by ISO3166 code (country)
	area["ISO3166-1"="CA"][admin_level=2]->.canada;

	// Find Edmonton area inside Canada (admin_level=6)
	area["name"="Edmonton"][admin_level=6](area.canada)->.edmonton_area;

	// Query ways inside Edmonton area
	(
		way["highway"]["area"!~"yes"]["highway"!~"motorway|motorway_link|raceway|construction|service"]["bicycle"!~"no"](area.edmonton_area);
		way["highway"="cycleway"]["bicycle"!~"no"](area.edmonton_area);
	);
	out body;
	>;
	out skel qt;
	`
	data.GetOSMData(query)
	*/

	fileName := filepath.Join("osm_bike_data.json")
	network, _ := data.BuildGraph(fileName)

	mux := http.NewServeMux()
	server.RegisterRoutes(mux, network)

	fileServer := http.FileServer(http.Dir("./web"))
	mux.Handle("/", fileServer)

	addr := ":8080"
	log.Printf("Starting server on %s\n", addr)
	err := http.ListenAndServe(addr, mux)
	if err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}