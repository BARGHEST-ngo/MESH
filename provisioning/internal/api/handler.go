package api

import (
	"crypto/sha256"
	"crypto/subtle"
	"net/http"
)

func NewRouter(apiKey string) http.Handler {
	keyHash := sha256.Sum256([]byte(apiKey))

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", handleHealth)
	mux.HandleFunc("POST /deployment", handleDeployment)

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
