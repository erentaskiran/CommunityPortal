package api

import (
	"database/sql"
	"net/http"

	"github.com/erentaskiran/project123123123/internal/middleware"
	"github.com/gorilla/mux"
)

type Router struct {
	db *sql.DB
}

func NewRouter(db *sql.DB) *Router {
	return &Router{
		db: db,
	}
}

func (r *Router) NewRouter() *mux.Router {

	router := mux.NewRouter()

	router.Use(middleware.CorsMiddleware)

	router.HandleFunc("/healthcheck", HealthCheck).Methods(http.MethodGet, http.MethodOptions)

	return router
}
