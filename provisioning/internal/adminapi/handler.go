package adminapi

import (
	"crypto/sha256"
	"crypto/subtle"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/BARGHEST-ngo/MESH/provisioning/internal/state"
)

type createKeyRequest struct {
	OwnerID       string `json:"owner_id"`
	Label         string `json:"label"`
	MaxConcurrent int    `json:"max_concurrent"`
	TTLHours      int    `json:"ttl_hours"` // 0 - No Expiry
}

type createKeyResponse struct {
	ID        string     `json:"id"`
	Key       string     `json:"key"`
	OwnerID   string     `json:"owner_id"`
	Label     string     `json:"label"`
	ExpiresAt *time.Time `json:"expires_at"`
}

type updateKeyRequest struct {
	Label         *string    `json:"label"`
	MaxConcurrent *int       `json:"max_concurrent"`
	ExpiresAt     *time.Time `json:"expires_at"`
	ClearExpiry   bool       `json:"clear_expiry"`
}

type getKeysResponse struct {
	Keys map[string]keyResponse `json:"keys"`
}

type keyResponse struct {
	ID            string     `json:"id"`
	OwnerID       string     `json:"owner_id"`
	Label         string     `json:"label"`
	CreatedAt     time.Time  `json:"created_at"`
	ExpiresAt     *time.Time `json:"expires_at"`
	MaxConcurrent int        `json:"max_concurrent"`
	Revoked       bool       `json:"revoked"`
}

type handler struct {
	keys *state.KeyStore
}

func NewAdminRouter(keys *state.KeyStore, adminToken string) http.Handler {
	h := &handler{
		keys: keys,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /keys", h.handleGetKeys)
	mux.HandleFunc("POST /keys", h.handlePostKey)
	mux.HandleFunc("DELETE /keys/{key_id}", h.handleDeleteKey)
	mux.HandleFunc("PATCH /keys/{key_id}", h.handlePatchKey)

	return authRequest(adminToken, mux)
}

func authRequest(adminToken string, next http.Handler) http.Handler {
	expected := sha256.Sum256([]byte(adminToken))
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

func (h *handler) handleGetKeys(w http.ResponseWriter, r *http.Request) {
	keys := h.keys.List()
	respKeys := make(map[string]keyResponse, len(keys))
	for _, k := range keys {
		respKeys[k.ID] = keyResponse{
			ID:            k.ID,
			OwnerID:       k.OwnerID,
			Label:         k.Label,
			CreatedAt:     k.CreatedAt,
			ExpiresAt:     k.ExpiresAt,
			MaxConcurrent: k.MaxConcurrent,
			Revoked:       k.Revoked,
		}
	}
	response := getKeysResponse{
		Keys: respKeys,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (h *handler) handlePostKey(w http.ResponseWriter, r *http.Request) {
	var req createKeyRequest
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.Label == "" || req.OwnerID == "" || req.MaxConcurrent < 0 {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	var ttl *time.Duration
	if req.TTLHours > 0 {
		dur := time.Hour * time.Duration(req.TTLHours)
		ttl = &dur
	}
	key, plaintext, err := h.keys.Create(req.OwnerID, req.Label, req.MaxConcurrent, ttl)
	if err != nil {
		http.Error(w, "failed to create key", http.StatusInternalServerError)
		return
	}

	response := createKeyResponse{
		ID:        key.ID,
		Key:       plaintext,
		OwnerID:   key.OwnerID,
		Label:     key.Label,
		ExpiresAt: key.ExpiresAt,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (h *handler) handleDeleteKey(w http.ResponseWriter, r *http.Request) {
	keyID := r.PathValue("key_id")
	err := h.keys.Revoke(keyID)
	if err != nil {
		if errors.Is(err, state.ErrNotFound) {
			http.Error(w, "unknown key_id", http.StatusNotFound)
			return
		}
		http.Error(w, fmt.Sprintf("failed to revoke key: %s", keyID), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// TODO: If maxConcurrent is less than current number of deployments - what to do?
// Don't kill existing deployments and just block new requests?
// Killing existing deployments requires choosing one and communicating that to the end-user
func (h *handler) handlePatchKey(w http.ResponseWriter, r *http.Request) {
	var req updateKeyRequest
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	keyID := r.PathValue("key_id")
	if req.MaxConcurrent != nil && *req.MaxConcurrent < 0 {
		http.Error(w, "max_concurrent must be >= 0", http.StatusBadRequest)
		return
	}

	if req.Label != nil && *req.Label == "" {
		http.Error(w, "label cannot be empty string", http.StatusBadRequest)
		return
	}

	// defaults to state.ExpiryNoChange
	var expiryUpdate state.ExpiryUpdate
	if req.ClearExpiry {
		expiryUpdate = state.ExpiryUpdate{Op: state.ExpiryClear}
	} else if req.ExpiresAt != nil {
		expiryUpdate = state.ExpiryUpdate{Op: state.ExpirySet, Value: req.ExpiresAt}
	}

	err := h.keys.Update(keyID, req.Label, req.MaxConcurrent, expiryUpdate)
	if err != nil {
		if errors.Is(err, state.ErrNotFound) {
			http.Error(w, "unknown key_id", http.StatusNotFound)
			return
		}
		http.Error(w, fmt.Sprintf("failed to update key: %s", keyID), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
