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
	"github.com/BARGHEST-ngo/MESH/provisioning/internal/reaper"
	"github.com/BARGHEST-ngo/MESH/provisioning/internal/state"
	"github.com/BARGHEST-ngo/MESH/provisioning/internal/workerclient"
	"golang.org/x/time/rate"
)

func main() {
	portMin := getEnvInt("FRPS_PORT_MIN")
	portMax := getEnvInt("FRPS_PORT_MAX")
	defaultTTLHours := getEnvInt("DEFAULT_TTL_HOURS")
	defaultTTL := time.Duration(defaultTTLHours) * time.Hour

	dataPath := getEnv("HOST_DATA_PATH")
	workerAddr := getEnv("WORKER_ADDR")
	workerToken := getEnv("WORKER_TOKEN")

	registry, err := state.New(filepath.Join(dataPath, "state.json"), portMin, portMax, defaultTTL)
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

	reaperCtx, reaperCancel := context.WithCancel(context.Background())
	defer reaperCancel()
	go reaper.Run(reaperCtx, registry, containerSvc, time.Minute*time.Duration(30))

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

func getEnv(variable string) string {
	str := os.Getenv(variable)
	if str == "" {
		log.Fatalf("'%s' must be set", variable)
	}
	return str
}

func getEnvInt(variable string) int {
	str := getEnv(variable)
	varInt, err := strconv.Atoi(str)
	if err != nil {
		log.Fatalf("failed to parse '%s'", variable)
	}
	return varInt
}
