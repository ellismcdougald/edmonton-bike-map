package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	"github.com/ellismcdougald/edmonton-bike-map/internal/domain/network"
	"github.com/ellismcdougald/edmonton-bike-map/internal/handler"
	"github.com/ellismcdougald/edmonton-bike-map/internal/repository/sqlrepo"
	internalserver "github.com/ellismcdougald/edmonton-bike-map/internal/server"
	"github.com/ellismcdougald/edmonton-bike-map/internal/service"
	"github.com/ellismcdougald/edmonton-bike-map/internal/token"
)

func main() {
	log.Printf("testing")
	_ = godotenv.Load()

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		user := os.Getenv("POSTGRES_USER")
		password := os.Getenv("POSTGRES_PASSWORD")
		dbname := os.Getenv("POSTGRES_DB")
		port := os.Getenv("POSTGRES_PORT")
		if port == "" {
			port = "5432"
		}
		dbURL = fmt.Sprintf("postgres://%s:%s@localhost:%s/%s?sslmode=disable", user, password, port, dbname)
	}

	log.Printf("connecting to db: %s", dbURL)
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Could not connect to db: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("Error closing db: %v", err)
		}
	}()

	if err := db.Ping(); err != nil {
		log.Fatalf("Could not ping db: %v", err)
	}

	userRepo := sqlrepo.NewSQLUserRepository(db)
	nodeRepo := sqlrepo.NewSQLNodeRepository(db)
	wayRepo := sqlrepo.NewSQLWayRepository(db)
	reviewRepo := sqlrepo.NewSQLReviewRepository(db)

	userSvc := service.NewUserService(userRepo)
	nodeSvc := service.NewNodeService(nodeRepo)
	waySvc := service.NewWayService(wayRepo)
	reviewSvc := service.NewReviewService(reviewRepo)

	netw, err := network.BuildNetwork(*nodeSvc, *waySvc, *reviewSvc)
	if err != nil {
		log.Fatalf("Error building network: %v", err)
	}

	routeSvc := service.NewRouteService(netw)

	jwtSecret := []byte(os.Getenv("JWT_SECRET"))
	tp := token.NewJWTProvider(jwtSecret)

	// handlers
	authHandler := handler.NewAuthHandler(userSvc, tp)
	userHandler := handler.NewUserHandler(userSvc)
	wayHandler := handler.NewWayHandler(nodeSvc, waySvc)
	reviewHandler := handler.NewReviewHandler(reviewSvc)
	routeHandler := handler.NewRouteHandler(routeSvc)

	handlers := handler.NewHandlers(authHandler, userHandler, wayHandler, reviewHandler, routeHandler)

	mux := http.NewServeMux()
	internalserver.RegisterRoutes(mux, *handlers)

	fileServer := http.FileServer(http.Dir("./web"))
	fileHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store")
		fileServer.ServeHTTP(w, r)
	})
	mux.Handle("/", fileHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	addr := ":" + port
	log.Printf("Starting new_server on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
