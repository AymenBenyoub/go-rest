package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"rest/server"
	"time"
)

func main() {

	srv := &http.Server{
		Addr:         ":8080",
		Handler:      server.NewRouter(),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Println("Starting server on :8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe failed: %v", err)
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt) // capture SIGINT
	<-quit                            // BLOCK until signal received

	log.Println("Shutting down server...")

	// Context to give current requests time to finish (max 10s)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited cleanly")
}
