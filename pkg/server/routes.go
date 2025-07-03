package server

import (
	"net/http"

	"github.com/ellismcdougald/edmonton-bike-map/pkg/model"
)

func RegisterRoutes(mux *http.ServeMux, network *model.Graph) {
	mux.HandleFunc("/api/route", func(writer http.ResponseWriter, request *http.Request) {
		handleRoute(writer, request, network)
	})
}