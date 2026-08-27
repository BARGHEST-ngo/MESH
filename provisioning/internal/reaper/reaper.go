package reaper

import (
	"context"
	"log/slog"
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
					slog.Error("failed to stop container", "err", err)
					continue
				}

				if err := registry.Release(d.Slug); err != nil {
					slog.Error("failed to release port", "err", err)
					continue
				}
				slog.Info("stopped and released port", "slug", d.Slug, "port", d.FrpsPort)
			}
		}
	}
}
