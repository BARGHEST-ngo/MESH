package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/BARGHEST-ngo/MESH/provisioning/internal/docker"
	"github.com/BARGHEST-ngo/MESH/provisioning/internal/env"
	"github.com/BARGHEST-ngo/MESH/provisioning/internal/workerapi"
)

func main() {
	frpsImage := env.GetEnv("FRPS_IMAGE")
	workerToken := env.GetEnv("WORKER_TOKEN")
	// Empty string defaults to "0.0.0.0"
	frpsBindAddr := os.Getenv("FRPS_BIND_ADDR")

	if err := docker.PullImage(frpsImage); err != nil {
		env.Fatal("failed to pull frps image", "err", err)
	}

	runner := docker.Manager{
		FrpsBindAddr: frpsBindAddr,
		FrpsImage:    frpsImage,
	}

	srv := &http.Server{
		Addr:         ":8081",
		Handler:      workerapi.NewWorkerRouter(runner, workerToken),
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	go func() {
		slog.Info("worker listening on :8081")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			env.Fatal("listen", "err", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("shutdown:", "err", err)
	}
}
