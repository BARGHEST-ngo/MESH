package admin

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/BARGHEST-ngo/MESH/provisioning/internal/state"
)

func newTestKeyStore(t *testing.T) *state.KeyStore {
	t.Helper()
	ks, err := state.NewKeyStore(filepath.Join(t.TempDir(), "keys.json"))
	if err != nil {
		t.Fatal(err)
	}
	return ks
}

func newTestAdminRouter(t *testing.T) http.Handler {
	t.Helper()
	return NewAdminRouter(newTestKeyStore(t))
}

func TestCreateKey(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		const ownerId = "test-owner"
		const label = "test-label"
		requestData := createKeyRequest{
			OwnerID:       ownerId,
			Label:         label,
			MaxConcurrent: 0,
			TTLHours:      0,
		}
		b, err := json.Marshal(&requestData)
		if err != nil {
			t.Fatal(err)
		}

		req := httptest.NewRequest(http.MethodPost, "/keys", bytes.NewBuffer(b))
		w := httptest.NewRecorder()
		newTestAdminRouter(t).ServeHTTP(w, req)
		if w.Code != http.StatusCreated {
			t.Errorf("expected %d, got %d", http.StatusOK, w.Code)
		}

		var response createKeyResponse
		if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
			t.Fatal(err)
		}

		if response.ID == "" || response.Key == "" || response.Label != label || response.OwnerID != ownerId || response.ExpiresAt != nil {
			t.Error("invalid response")
		}
	})

	t.Run("invalid-method", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/keys", nil)
		w := httptest.NewRecorder()
		newTestAdminRouter(t).ServeHTTP(w, req)
		if w.Code != http.StatusMethodNotAllowed {
			t.Errorf("expected %d, got %d", http.StatusMethodNotAllowed, w.Code)
		}
	})

	t.Run("invalid-data", func(t *testing.T) {
		b, err := json.Marshal(&createKeyRequest{})
		if err != nil {
			t.Fatal(err)
		}

		req := httptest.NewRequest(http.MethodPost, "/keys", bytes.NewBuffer(b))
		w := httptest.NewRecorder()
		newTestAdminRouter(t).ServeHTTP(w, req)
		if w.Code != http.StatusBadRequest {
			t.Errorf("expected %d, got %d", http.StatusBadRequest, w.Code)
		}
	})

	t.Run("invalid-data", func(t *testing.T) {
		b, err := json.Marshal(&createKeyRequest{
			OwnerID: "",
			Label:   "label",
		})
		if err != nil {
			t.Fatal(err)
		}

		req := httptest.NewRequest(http.MethodPost, "/keys", bytes.NewBuffer(b))
		w := httptest.NewRecorder()
		newTestAdminRouter(t).ServeHTTP(w, req)
		if w.Code != http.StatusBadRequest {
			t.Errorf("expected %d, got %d", http.StatusBadRequest, w.Code)
		}
	})

	t.Run("invalid-data", func(t *testing.T) {
		b, err := json.Marshal(&createKeyRequest{
			OwnerID: "owner",
			Label:   "",
		})
		if err != nil {
			t.Fatal(err)
		}

		req := httptest.NewRequest(http.MethodPost, "/keys", bytes.NewBuffer(b))
		w := httptest.NewRecorder()
		newTestAdminRouter(t).ServeHTTP(w, req)
		if w.Code != http.StatusBadRequest {
			t.Errorf("expected %d, got %d", http.StatusBadRequest, w.Code)
		}
	})

	t.Run("invalid-data", func(t *testing.T) {
		b, err := json.Marshal(&createKeyRequest{
			OwnerID:       "owner",
			Label:         "label",
			MaxConcurrent: -1,
		})
		if err != nil {
			t.Fatal(err)
		}

		req := httptest.NewRequest(http.MethodPost, "/keys", bytes.NewBuffer(b))
		w := httptest.NewRecorder()
		newTestAdminRouter(t).ServeHTTP(w, req)
		if w.Code != http.StatusBadRequest {
			t.Errorf("expected %d, got %d", http.StatusBadRequest, w.Code)
		}
	})
}
