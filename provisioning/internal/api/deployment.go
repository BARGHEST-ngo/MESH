package api

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
)

// Provision a new frp container
// Generate subdomain slug, token & return to client
func (h *handler) handlePostDeployment(w http.ResponseWriter, r *http.Request) {
	slug, err := generateSlug()
	if err != nil {
		http.Error(w, "failed to generate slug", http.StatusInternalServerError)
		return
	}

	token, err := generateToken()
	if err != nil {
		http.Error(w, "failed to generate token", http.StatusInternalServerError)
		return
	}

	port, err := h.registry.AllocatePort(slug, token)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to allocate port: %v", err), http.StatusInternalServerError)
		return
	}

	// start container

	response := &DeploymentResponse{
		Slug:     slug,
		Token:    token,
		FrpsPort: port,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func generateSlug() (string, error) {
	b := make([]byte, 5)
	_, err := rand.Read(b)
	return hex.EncodeToString(b), err
}

func generateToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	return base64.RawURLEncoding.EncodeToString(b), err
}

func (h *handler) handleDeleteDeployment(w http.ResponseWriter, r *http.Request) {}
