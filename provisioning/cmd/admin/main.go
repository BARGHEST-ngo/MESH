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

	"github.com/BARGHEST-ngo/MESH/provisioning/internal/adminapi"
	"github.com/BARGHEST-ngo/MESH/provisioning/internal/env"
	"github.com/BARGHEST-ngo/MESH/provisioning/internal/state"
)

func main() {
	dataPath := env.GetEnv("HOST_DATA_PATH")
	adminToken := env.GetEnv("ADMIN_TOKEN")

	keyStore, err := state.NewKeyStore(filepath.Join(dataPath, "keys.json"))
	if err != nil {
		env.Fatal("failed to initialise key store")
	}

	srv := &http.Server{
		Addr:         ":9090",
		Handler:      adminapi.NewAdminRouter(keyStore, adminToken),
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	go func() {
		slog.Info("admin listening on :9090")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			env.Fatal("listen: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("shutdown: %v", "err", err)
	}
}
