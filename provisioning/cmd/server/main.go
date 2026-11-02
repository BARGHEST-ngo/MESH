package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/BARGHEST-ngo/MESH/provisioning/internal/api"
	"github.com/BARGHEST-ngo/MESH/provisioning/internal/env"
	"github.com/BARGHEST-ngo/MESH/provisioning/internal/reaper"
	"github.com/BARGHEST-ngo/MESH/provisioning/internal/state"
	"github.com/BARGHEST-ngo/MESH/provisioning/internal/workerclient"
	"golang.org/x/time/rate"
)

func main() {
	portMin := env.GetEnvInt("FRPS_PORT_MIN")
	portMax := env.GetEnvInt("FRPS_PORT_MAX")
	defaultTTLHours := env.GetEnvInt("DEFAULT_TTL_HOURS")
	defaultTTL := time.Duration(defaultTTLHours) * time.Hour

	dataPath := env.GetEnv("HOST_DATA_PATH")
	workerAddr := env.GetEnv("WORKER_ADDR")
	workerToken := env.GetEnv("WORKER_TOKEN")

	registry, err := state.New(filepath.Join(dataPath, "state.json"), portMin, portMax, defaultTTL)
	if err != nil {
		env.Fatal("failed to initialise port registry")
	}

	keyStore, err := state.NewKeyStore(filepath.Join(dataPath, "keys.json"))
	if err != nil {
		env.Fatal("failed to initialise key store")
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

	reaperCtx, reaperCancel := context.WithCancel(context.Background())
	defer reaperCancel()
	go reaper.Run(reaperCtx, registry, containerSvc, time.Minute*time.Duration(30))

	go func() {
		slog.Info("provisioner listening on :8080")
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
