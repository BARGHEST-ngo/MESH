package api

import (
	"context"
	"crypto/sha256"
	"net/http"

	"github.com/BARGHEST-ngo/MESH/provisioning/internal/state"
)

type handler struct {
	registry *state.Registry
	service  ContainerService
}

type ContainerService interface {
	Start(d state.Deployment) error
	Stop(slug string) error
}

type contextKey string

const apiKeyContextKey contextKey = "apikey"

func NewRouter(keys *state.KeyStore, registry *state.Registry, containerSvc ContainerService) http.Handler {
	h := &handler{
		registry: registry,
		service:  containerSvc,
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

		bearer := r.Header.Get("Authorization")
		if len(bearer) < 8 || bearer[:7] != "Bearer " {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		token := sha256.Sum256([]byte(bearer[7:]))
		key, ok := keys.Lookup(token)
		if !ok {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), apiKeyContextKey, key)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
