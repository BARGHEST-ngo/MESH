package api

import (
	"crypto/sha256"
	"crypto/subtle"
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

func NewRouter(apiKey string, registry *state.Registry, frpsImage, meshDomain string, opts ...Option) http.Handler {
	keyHash := sha256.Sum256([]byte(apiKey))
	h := &handler{
		registry: registry,
		service: docker.Manager{
			FrpsImage:  frpsImage,
			MeshDomain: meshDomain,
		},
	}
	for _, opt := range opts {
		opt(h)
	}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", handleHealth)
	mux.HandleFunc("POST /deployment", h.handlePostDeployment)
	mux.HandleFunc("DELETE /deployment/{slug}", h.handleDeleteDeployment)

	return authRequest(keyHash, mux)
}

func authRequest(keyHash [32]byte, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/health" {
			next.ServeHTTP(w, r)
			return
		}

		bearer := r.Header.Get("Authorization")
		if len(bearer) < 8 || bearer[:7] != "Bearer " {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		incoming := sha256.Sum256([]byte(bearer[7:]))
		if subtle.ConstantTimeCompare(incoming[:], keyHash[:]) != 1 {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
