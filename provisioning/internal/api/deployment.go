package api

import (
	"encoding/json"
	"net/http"
)

// Provision a new frp container
// Generate subdomain slug, token & return to client
func handleDeployment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"deployment": "ok"})
}
