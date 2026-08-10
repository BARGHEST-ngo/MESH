package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/BARGHEST-ngo/MESH/provisioning/internal/docker"
	"github.com/BARGHEST-ngo/MESH/provisioning/internal/workerapi"
)

func main() {
	if os.Getenv("HOST_DATA_PATH") == "" {
		log.Fatal("HOST_DATA_PATH must be set")
	}

	frpsImage := os.Getenv("FRPS_IMAGE")
	if frpsImage == "" {
		log.Fatal("FRPS_IMAGE must be set")
	}

	workerToken := os.Getenv("WORKER_TOKEN")
	if workerToken == "" {
		log.Fatal("WORKER_TOKEN must be set")
	}

	if err := docker.PullImage(frpsImage); err != nil {
		log.Fatalf("failed to pull frps image: %v", err)
	}

	runner := docker.Manager{FrpsImage: frpsImage}

	srv := &http.Server{
		Addr:         ":8081",
		Handler:      workerapi.NewWorkerRouter(runner, workerToken),
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	go func() {
		log.Printf("worker listening on :8081")
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
