package middleware

import (
	"context"
	"net/http"
	"time"
)

// TimeoutMiddleware sets a timeout for the request context (pass it to handlers & repositories)

func TimeoutMiddleware(timeout int) Middleware {
	return func (next http.Handler) http.Handler {
      return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), time.Second*time.Duration(timeout))
		defer cancel()
		next.ServeHTTP(w, r.WithContext(ctx))
	  })
		

	}
}
