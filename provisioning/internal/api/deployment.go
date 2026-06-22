package api

import (
	"crypto/rand"
	"encoding/json"
	"net/http"
)

// Provision a new frp container
// Generate subdomain slug, token & return to client
func handlePostDeployment(w http.ResponseWriter, r *http.Request) {
	slug, err := generateSlug()
	if err != nil {
		// write error
		return
	}

	response := &DeploymentResponse{
		Slug: slug,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func generateSlug() (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, 10)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	for i := range b {
		b[i] = charset[int(b[i])%len(charset)]
	}
	return string(b), nil
}

func handleDeleteDeployment(w http.ResponseWriter, r *http.Request) {}
