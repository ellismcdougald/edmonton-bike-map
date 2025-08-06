package handlers

import (
	"database/sql"
	"net/http"

	"github.com/ellismcdougald/edmonton-bike-map/pkg/model"
)

type APIHandlers interface {
	HandleLogin() http.HandlerFunc
	HandleSignup() http.HandlerFunc
}

type RealHandlers struct {
	UserService model.UserService
	DB *sql.DB
	Network *model.Graph
}