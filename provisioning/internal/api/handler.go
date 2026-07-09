package api

import (
	"context"
	"crypto/sha256"
	"fmt"
	"net/http"

	"github.com/BARGHEST-ngo/MESH/provisioning/internal/docker"
	"github.com/BARGHEST-ngo/MESH/provisioning/internal/state"
)

type handler struct {
	registry *state.Registry
	service  ContainerService
}

type Option func(*handler)

func WithContainerService(svc ContainerService) Option {
	return func(h *handler) { h.service = svc }
}

type ContainerService interface {
	Start(d state.Deployment) error
	Stop(slug string) error
}

type contextKey string

const apiKeyContextKey contextKey = "ownerID"

func NewRouter(keys *state.KeyStore, registry *state.Registry, frpsImage string, opts ...Option) http.Handler {
	h := &handler{
		registry: registry,
		service: docker.Manager{
			FrpsImage: frpsImage,
		},
	}
	for _, opt := range opts {
		opt(h)
	}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", handleHealth)
	mux.HandleFunc("POST /deployment", h.handlePostDeployment)
	mux.HandleFunc("DELETE /deployment/{slug}", h.handleDeleteDeployment)

	return authRequest(keys, mux)
}

func authRequest(keys *state.KeyStore, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/health" {
			next.ServeHTTP(w, r)
			return
		}

		incoming, err := getAuthToken(r)
		if err != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		key, ok := keys.Lookup(incoming)
		if !ok {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), apiKeyContextKey, key)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getAuthToken(r *http.Request) ([32]byte, error) {
	bearer := r.Header.Get("Authorization")
	if len(bearer) < 8 || bearer[:7] != "Bearer " {
		return [32]byte{}, fmt.Errorf("invalid auth token")
	}

	token := sha256.Sum256([]byte(bearer[7:]))
	return token, nil
}
