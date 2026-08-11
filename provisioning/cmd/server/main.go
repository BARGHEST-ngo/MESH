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
	"github.com/BARGHEST-ngo/MESH/provisioning/internal/workerclient"
	"golang.org/x/time/rate"
)

func main() {
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

	workerAddr := os.Getenv("WORKER_ADDR")
	if workerAddr == "" {
		log.Fatal("WORKER_ADDR must be set")
	}

	workerToken := os.Getenv("WORKER_TOKEN")
	if workerToken == "" {
		log.Fatal("WORKER_TOKEN must be set")
	}

	registry, err := state.New(filepath.Join(dataPath, "state.json"), portMinInt, portMaxInt)
	if err != nil {
		log.Fatal("failed to initialise port registry")
	}

	keyStore, err := state.NewKeyStore(filepath.Join(dataPath, "keys.json"))
	if err != nil {
		log.Fatal("failed to initialise key store")
	}

	deploymentRateLimit := rate.Every(10 * time.Second)
	deploymentRateBurst := 3
	rateLimiter := api.NewLimiter(deploymentRateLimit, deploymentRateBurst)
	containerSvc := workerclient.New(workerAddr, workerToken)
	srv := &http.Server{
		Addr:         ":8080",
		Handler:      api.NewRouter(keyStore, registry, containerSvc, rateLimiter),
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
