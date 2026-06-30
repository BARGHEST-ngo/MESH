package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"syscall"
	"time"

	"github.com/BARGHEST-ngo/MESH/provisioning/internal/api"
	"github.com/BARGHEST-ngo/MESH/provisioning/internal/state"
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

	portMin := os.Getenv("FRPS_PORT_MIN")
	if portMin == "" {
		log.Fatal("FRPS_PORT_MIN must be set")
	}

	portMinInt, err := strconv.Atoi(portMin)
	if err != nil {
		log.Fatal("failed to parse FRPS_PORT_MIN")
	}

	portMax := os.Getenv("FRPS_PORT_MAX")
	if portMax == "" {
		log.Fatal("FRPS_PORT_MAX must be set")
	}
	portMaxInt, err := strconv.Atoi(portMax)
	if err != nil {
		log.Fatal("failed to parse FRPS_PORT_MAX")
	}

	dataPath := os.Getenv("HOST_DATA_PATH")
	if dataPath == "" {
		log.Fatal("HOST_DATA_PATH must be set")
	}

	frpsImage := os.Getenv("FRPS_IMAGE")
	if frpsImage == "" {
		log.Fatal("FRPS_IMAGE must be set")
	}

	registry, err := state.New(filepath.Join(dataPath, "state.json"), portMinInt, portMaxInt)
	if err != nil {
		log.Fatal("failed to initialise port registry")
	}

	srv := &http.Server{
		Addr:         ":8080",
		Handler:      api.NewRouter(apiKey, registry, frpsImage),
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
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
