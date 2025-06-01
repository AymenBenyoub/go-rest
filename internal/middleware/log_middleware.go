package middleware

import (
	"log"
	"net/http"
	"time"
)

// log the details of every request
func RequestLogMiddleware() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			next.ServeHTTP(w, r) // call next middleware or final handler

			duration := time.Since(start)

			log.Printf("Request %s %s from %s completed in %v", r.Method, r.URL.Path, r.RemoteAddr, duration)

		})
	}
}

