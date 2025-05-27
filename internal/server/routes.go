package server

import (
	"database/sql"
	"net/http"
)

func NewRouter(db *sql.DB) http.Handler {

	router := http.NewServeMux()

	router.HandleFunc("GET /posts")

	return router
}
