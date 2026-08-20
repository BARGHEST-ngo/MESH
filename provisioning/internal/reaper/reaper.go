package reaper

import (
	"context"
	"log"
	"time"

	"github.com/BARGHEST-ngo/MESH/provisioning/internal/api"
	"github.com/BARGHEST-ngo/MESH/provisioning/internal/state"
)

func Run(ctx context.Context, registry *state.Registry, service api.ContainerService, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case t := <-ticker.C:
			expired := registry.Expired(t)
			for _, d := range expired {
				if err := service.Stop(d.Slug); err != nil {
					log.Printf("failed to stop container: %v", err)
					continue
				}

				if err := registry.Release(d.Slug); err != nil {
					log.Printf("failed to release port: %v", err)
					continue
				}
				log.Printf("%s stopped and released port %d", d.Slug, d.FrpsPort)
			}
		}
	}
}
