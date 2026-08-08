package api

import (
	"context"
	"crypto/sha256"
	"net/http"

	"github.com/BARGHEST-ngo/MESH/provisioning/internal/state"
)

type handler struct {
	registry    *state.Registry
	service     ContainerService
	rateLimiter *rateLimiter
}

type ContainerService interface {
	Start(d state.Deployment) error
	Stop(slug string) error
}

type contextKey string

const apiKeyContextKey contextKey = "apikey"

func NewRouter(keys *state.KeyStore, registry *state.Registry, containerSvc ContainerService, rateLimiter *rateLimiter) http.Handler {
	h := &handler{
		registry:    registry,
		service:     containerSvc,
		rateLimiter: rateLimiter,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", handleHealth)
	mux.HandleFunc("POST /deployment", rateLimit(h.rateLimiter, h.handlePostDeployment))
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

func rateLimit(l *rateLimiter, next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key, ok := r.Context().Value(apiKeyContextKey).(state.APIKey)
		if !ok {
			http.Error(w, "invalid context", http.StatusBadRequest)
			return
		}

		if !l.allow(key.ID) {
			http.Error(w, "too many requests", http.StatusTooManyRequests)
			return
		}

		next(w, r)
	})
}
