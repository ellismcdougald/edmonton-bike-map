package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/ellismcdougald/edmonton-bike-map/pkg/data"
	"github.com/ellismcdougald/edmonton-bike-map/pkg/server"
	"github.com/joho/godotenv"
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

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("POSTGRES_DB")
	port := os.Getenv("POSTGRES_PORT")

	connectionStr := fmt.Sprintf("postgres://%s:%s@localhost:%s/%s?sslmode=disable", user, password, port, dbname)

	db, err := sql.Open("postgres", connectionStr)
	if err != nil {
		log.Fatalf("Could not connect to db: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("Error closing db: %v", err)
		}
	}()

	fileName := filepath.Join("osm_bike_data.json")
	network, _ := data.BuildGraph(fileName)
	//allData, _ := data.GetAllGeoJsonData(fileName)

	mux := http.NewServeMux()
	server.RegisterRoutes(mux, network, db)

	fileServer := http.FileServer(http.Dir("./web"))
	fileHandler := http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Cache-Control", "no-store")
		fileServer.ServeHTTP(writer, request)
	})
	mux.Handle("/", fileHandler)

	addr := ":8080"
	log.Printf("Starting server on %s\n", addr)
	err = http.ListenAndServe(addr, mux)
	if err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
