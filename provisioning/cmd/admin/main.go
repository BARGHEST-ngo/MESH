package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/BARGHEST-ngo/MESH/provisioning/internal/adminapi"
	"github.com/BARGHEST-ngo/MESH/provisioning/internal/state"
)

func main() {
	dataPath := os.Getenv("HOST_DATA_PATH")
	if dataPath == "" {
		log.Fatal("HOST_DATA_PATH must be set")
	}

	adminToken := os.Getenv("ADMIN_TOKEN")
	if adminToken == "" {
		log.Fatal("ADMIN_TOKEN must be set")
	}

	keyStore, err := state.NewKeyStore(filepath.Join(dataPath, "keys.json"))
	if err != nil {
		log.Fatal("failed to initialise key store")
	}

	srv := &http.Server{
		Addr:         ":9090",
		Handler:      adminapi.NewAdminRouter(keyStore, adminToken),
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	go func() {
		log.Printf("admin listening on :9090")
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
