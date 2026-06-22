package api

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
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

	token, err := generateToken()
	if err != nil {
		// write error
		return
	}

	response := &DeploymentResponse{
		Slug:  slug,
		Token: token,
		
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

func handleDeleteDeployment(w http.ResponseWriter, r *http.Request) {}
