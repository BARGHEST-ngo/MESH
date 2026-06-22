package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/BARGHEST-ngo/MESH/provisioning/internal/api"
)

// This is an intentionally light-weight and basic HTTP server
// Don't want to over-engineer at this stage, just prove it works
func main() {
	// PROVISIONING_API_KEY during early dev will be a single shared secret
	// Internal use only during development
	// Moving to per-customer keys as we get closer to prod
	apiKey := os.Getenv("PROVISIONING_API_KEY")
	if apiKey == "" {
		log.Fatal("PROVISIONING_API_KEY must be set")
	}

	srv := &http.Server{
		Addr:         ":8080",
		Handler:      api.NewRouter(apiKey),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go func() {
		log.Printf("provisioner listening on :8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("shutdown: %v", err)
	}
}
