package workerapi

import (
	"crypto/sha256"
	"crypto/subtle"
	"encoding/json"
	"net/http"
	"regexp"

	"github.com/BARGHEST-ngo/MESH/provisioning/internal/state"
)

type ContainerRunner interface {
	Start(d state.Deployment) error
	Stop(slug string) error
}

type handler struct {
	runner ContainerRunner
}

func NewWorkerRouter(runner ContainerRunner, workerToken string) http.Handler {
	h := &handler{runner: runner}

	mux := http.NewServeMux()
	mux.HandleFunc("POST /containers", h.handleStart)
	mux.HandleFunc("DELETE /containers/{slug}", h.handleStop)

	return authRequest(workerToken, mux)
}

func authRequest(workerToken string, next http.Handler) http.Handler {
	expected := sha256.Sum256([]byte(workerToken))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bearer := r.Header.Get("Authorization")
		if len(bearer) < 8 || bearer[:7] != "Bearer " {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		token := sha256.Sum256([]byte(bearer[7:]))
		if subtle.ConstantTimeCompare(token[:], expected[:]) != 1 {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (h *handler) handleStart(w http.ResponseWriter, r *http.Request) {
	var d state.Deployment
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&d); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.runner.Start(d); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *handler) handleStop(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")
	if !confirmValidSlug(slug) {
		http.Error(w, "invalid slug", http.StatusBadRequest)
		return
	}

	if err := h.runner.Stop(slug); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// mirrors confirmValidSlug in internal/api/deployment.go
func confirmValidSlug(slug string) bool {
	// 10 lowercase hex characters
	slugPattern := regexp.MustCompile(`^[0-9a-f]{10}$`)
	return slugPattern.MatchString(slug)
}
