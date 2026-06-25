package api

import (
	"crypto/sha256"
	"crypto/subtle"
	"net/http"

	"github.com/BARGHEST-ngo/MESH/provisioning/internal/docker"
	"github.com/BARGHEST-ngo/MESH/provisioning/internal/state"
)

type handler struct {
	registry       *state.Registry
	startContainer func(*state.Registry, string) error
}

type Option func(*handler)

func WithStartContainer(fn func(*state.Registry, string) error) Option {
	return func(h *handler) { h.startContainer = fn }
}

func NewRouter(apiKey string, registry *state.Registry, opts ...Option) http.Handler {
	keyHash := sha256.Sum256([]byte(apiKey))
	h := &handler{
		registry: registry,
		// Injecting this here for now so tests can complete
		// replace with interface and mocks for testing
		startContainer: docker.Start,
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
