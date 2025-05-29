package middleware

import (
	"log"
	"net/http"
	"time"
)

// log the details of every request

func LogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("Received request: %s %s from %s", r.Method, r.URL.Path, r.RemoteAddr)
		next.ServeHTTP(w, r)
		duration := time.Since(start)
		log.Printf("Request %s %s completed in %v", r.Method, r.URL.Path, duration)

	})
}
