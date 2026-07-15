package admin

import (
	"encoding/json"
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

type handler struct {
	keys *state.KeyStore
}

func NewAdminRouter(keys *state.KeyStore) http.Handler {
	h := &handler{
		keys: keys,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("POST /keys", h.handlePostKey)
	mux.HandleFunc("DELETE /keys/{key_id}", h.handleDeleteKey)
	mux.HandleFunc("PATCH /keys/{key_id}", h.handlePatchKey)

	return mux
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
	if keyID == "" {
		http.Error(w, "unknown key ID", http.StatusBadRequest)
		return
	}

	if err := h.keys.Revoke(keyID); err != nil {
		http.Error(w, fmt.Sprintf("failed to revoke key: %s", keyID), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *handler) handlePatchKey(w http.ResponseWriter, r *http.Request) {
	var req updateKeyRequest
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	keyID := r.PathValue("key_id")
	if keyID == "" {
		http.Error(w, "unknown key ID", http.StatusBadRequest)
		return
	}

	// defaults to state.ExpiryNoChange
	var expiryUpdate state.ExpiryUpdate
	if req.ClearExpiry {
		expiryUpdate = state.ExpiryUpdate{Op: state.ExpiryClear}
	} else if req.ExpiresAt != nil {
		expiryUpdate = state.ExpiryUpdate{Op: state.ExpirySet, Value: req.ExpiresAt}
	}

	if err := h.keys.Update(keyID, req.Label, req.MaxConcurrent, expiryUpdate); err != nil {
		http.Error(w, fmt.Sprintf("failed to revoke key: %s", keyID), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
