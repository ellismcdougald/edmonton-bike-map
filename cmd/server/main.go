package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/ellismcdougald/edmonton-bike-map/pkg/data"
	"github.com/ellismcdougald/edmonton-bike-map/pkg/model"
	"github.com/ellismcdougald/edmonton-bike-map/pkg/routing"
	"github.com/ellismcdougald/edmonton-bike-map/pkg/server"
	"github.com/ellismcdougald/edmonton-bike-map/pkg/server/handlers"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	user, password, dbname, port := getDBEnv()

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
	userStore := &model.DBUserStore{DB: db}
	nodeStore := &model.DBNodeStore{DB: db}
	wayStore := &model.DBWayStore{DB: db}
	reviewStore := &model.DBReviewStore{DB: db}
	router := &routing.RealRouter{}
	handlerFuncs := handlers.RealHandlers{
		UserService: userStore,
		NodeService: nodeStore,
		WayService: wayStore,
		ReviewService: reviewStore,
		DB:          db,
		Network:     network,
		Router: router,
	}
	server.RegisterRoutes(mux, network, db, &handlerFuncs)

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

func getDBEnv() (user, password, dbname, port string) {
    if os.Getenv("TESTING") == "true" {
        user = os.Getenv("POSTGRES_TEST_USER")
        password = os.Getenv("POSTGRES_TEST_PASSWORD")
        dbname = os.Getenv("POSTGRES_TEST_DB")
        port = os.Getenv("POSTGRES_TEST_PORT")
    } else {
        user = os.Getenv("POSTGRES_USER")
        password = os.Getenv("POSTGRES_PASSWORD")
        dbname = os.Getenv("POSTGRES_DB")
        port = os.Getenv("POSTGRES_PORT")
    }
    return
}
