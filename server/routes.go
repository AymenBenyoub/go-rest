package server

import "net/http"

func NewRouter() http.Handler {

	router := http.NewServeMux()

	router.HandleFunc("GET /", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Write([]byte("WAZZAP!!"))

	}))
  

	return router
}
