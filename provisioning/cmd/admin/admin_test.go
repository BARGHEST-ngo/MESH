package admin

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"
	"time"

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
			t.Errorf("expected %d, got %d", http.StatusCreated, w.Code)
		}

		var response createKeyResponse
		if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
			t.Fatal(err)
		}

		if response.ID == "" || response.Key == "" || response.Label != label || response.OwnerID != ownerId || response.ExpiresAt != nil {
			t.Error("invalid response")
		}
	})

	t.Run("valid-with-ttl", func(t *testing.T) {
		const ownerId = "test-owner"
		const label = "test-label"
		requestData := createKeyRequest{
			OwnerID:       ownerId,
			Label:         label,
			MaxConcurrent: 0,
			TTLHours:      1,
		}
		b, err := json.Marshal(&requestData)
		if err != nil {
			t.Fatal(err)
		}

		req := httptest.NewRequest(http.MethodPost, "/keys", bytes.NewBuffer(b))
		w := httptest.NewRecorder()
		newTestAdminRouter(t).ServeHTTP(w, req)
		if w.Code != http.StatusCreated {
			t.Errorf("expected %d, got %d", http.StatusCreated, w.Code)
		}

		var response createKeyResponse
		if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
			t.Fatal(err)
		}

		if response.ID == "" || response.Key == "" || response.Label != label || response.OwnerID != ownerId || response.ExpiresAt == nil {
			t.Error("invalid response")
		}
	})

	t.Run("empty-data", func(t *testing.T) {
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

	t.Run("malformed-json", func(t *testing.T) {
		b := []byte("not json")
		req := httptest.NewRequest(http.MethodPost, "/keys", bytes.NewBuffer(b))
		w := httptest.NewRecorder()
		newTestAdminRouter(t).ServeHTTP(w, req)
		if w.Code != http.StatusBadRequest {
			t.Errorf("expected %d, got %d", http.StatusBadRequest, w.Code)
		}
	})

	t.Run("missing-owner", func(t *testing.T) {
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

	t.Run("missing-label", func(t *testing.T) {
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

	t.Run("negative-maxConcurrent", func(t *testing.T) {
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

func TestDeleteKey(t *testing.T) {
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

		createReq := httptest.NewRequest(http.MethodPost, "/keys", bytes.NewBuffer(b))
		w := httptest.NewRecorder()
		router := newTestAdminRouter(t)
		router.ServeHTTP(w, createReq)
		if w.Code != http.StatusCreated {
			t.Errorf("expected %d, got %d", http.StatusCreated, w.Code)
		}

		var response createKeyResponse
		if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
			t.Fatal(err)
		}

		deleteReq := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/keys/%s", response.ID), nil)
		w = httptest.NewRecorder()
		router.ServeHTTP(w, deleteReq)
		if w.Code != http.StatusOK {
			t.Errorf("expected %d, got %d", http.StatusOK, w.Code)
		}
	})

	t.Run("unknown-key", func(t *testing.T) {
		deleteReq := httptest.NewRequest(http.MethodDelete, "/keys/unknown-id", nil)
		w := httptest.NewRecorder()
		newTestAdminRouter(t).ServeHTTP(w, deleteReq)
		if w.Code != http.StatusNotFound {
			t.Errorf("expected %d, got %d", http.StatusNotFound, w.Code)
		}
	})
}

func createPatchTestKey(t *testing.T, router http.Handler) createKeyResponse {
	t.Helper()
	requestData := createKeyRequest{
		OwnerID:       "test-owner",
		Label:         "test-label",
		MaxConcurrent: 0,
		TTLHours:      0,
	}
	b, err := json.Marshal(&requestData)
	if err != nil {
		t.Fatal(err)
	}

	createReq := httptest.NewRequest(http.MethodPost, "/keys", bytes.NewBuffer(b))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, createReq)
	if w.Code != http.StatusCreated {
		t.Fatal("expected success")
	}
	var response createKeyResponse
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatal(err)
	}
	return response
}

func TestPatchKey(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		router := newTestAdminRouter(t)
		response := createPatchTestKey(t, router)
		newLabel := "new-label"
		maxConcurrent := 1
		patchData := updateKeyRequest{
			Label:         &newLabel,
			MaxConcurrent: &maxConcurrent,
			ExpiresAt:     nil,
			ClearExpiry:   false,
		}
		b, err := json.Marshal(&patchData)
		if err != nil {
			t.Fatal(err)
		}
		patchReq := httptest.NewRequest(http.MethodPatch, fmt.Sprintf("/keys/%s", response.ID), bytes.NewBuffer(b))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, patchReq)
		if w.Code != http.StatusOK {
			t.Errorf("expected %d, got %d", http.StatusOK, w.Code)
		}
	})

	t.Run("valid-no-change", func(t *testing.T) {
		router := newTestAdminRouter(t)
		response := createPatchTestKey(t, router)

		patchData := updateKeyRequest{}
		b, err := json.Marshal(&patchData)
		if err != nil {
			t.Fatal(err)
		}
		patchReq := httptest.NewRequest(http.MethodPatch, fmt.Sprintf("/keys/%s", response.ID), bytes.NewBuffer(b))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, patchReq)
		if w.Code != http.StatusOK {
			t.Errorf("expected %d, got %d", http.StatusOK, w.Code)
		}
	})

	t.Run("clear-expiry", func(t *testing.T) {
		router := newTestAdminRouter(t)
		response := createPatchTestKey(t, router)
		patchData := updateKeyRequest{
			ClearExpiry: true,
		}
		b, err := json.Marshal(&patchData)
		if err != nil {
			t.Fatal(err)
		}
		patchReq := httptest.NewRequest(http.MethodPatch, fmt.Sprintf("/keys/%s", response.ID), bytes.NewBuffer(b))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, patchReq)
		if w.Code != http.StatusOK {
			t.Errorf("expected %d, got %d", http.StatusOK, w.Code)
		}
	})

	t.Run("set-expiry", func(t *testing.T) {
		router := newTestAdminRouter(t)
		response := createPatchTestKey(t, router)
		expiry := time.Now().Add(time.Hour)
		patchData := updateKeyRequest{
			ExpiresAt: &expiry,
		}
		b, err := json.Marshal(&patchData)
		if err != nil {
			t.Fatal(err)
		}
		patchReq := httptest.NewRequest(http.MethodPatch, fmt.Sprintf("/keys/%s", response.ID), bytes.NewBuffer(b))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, patchReq)
		if w.Code != http.StatusOK {
			t.Errorf("expected %d, got %d", http.StatusOK, w.Code)
		}
	})

	t.Run("invalid-label", func(t *testing.T) {
		router := newTestAdminRouter(t)
		response := createPatchTestKey(t, router)

		newLabel := ""
		patchData := updateKeyRequest{
			Label: &newLabel,
		}
		b, err := json.Marshal(&patchData)
		if err != nil {
			t.Fatal(err)
		}
		patchReq := httptest.NewRequest(http.MethodPatch, fmt.Sprintf("/keys/%s", response.ID), bytes.NewBuffer(b))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, patchReq)
		if w.Code != http.StatusBadRequest {
			t.Errorf("expected %d, got %d", http.StatusBadRequest, w.Code)
		}
	})

	t.Run("invalid-max-concurrent", func(t *testing.T) {
		router := newTestAdminRouter(t)
		response := createPatchTestKey(t, router)
		newMaxConcurrent := -1
		patchData := updateKeyRequest{
			MaxConcurrent: &newMaxConcurrent,
		}
		b, err := json.Marshal(&patchData)
		if err != nil {
			t.Fatal(err)
		}
		patchReq := httptest.NewRequest(http.MethodPatch, fmt.Sprintf("/keys/%s", response.ID), bytes.NewBuffer(b))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, patchReq)
		if w.Code != http.StatusBadRequest {
			t.Errorf("expected %d, got %d", http.StatusBadRequest, w.Code)
		}
	})

	t.Run("unknown-key", func(t *testing.T) {
		patchData := updateKeyRequest{}
		b, err := json.Marshal(&patchData)
		if err != nil {
			t.Fatal(err)
		}
		deleteReq := httptest.NewRequest(http.MethodPatch, "/keys/unknown-id", bytes.NewBuffer(b))
		w := httptest.NewRecorder()
		newTestAdminRouter(t).ServeHTTP(w, deleteReq)
		if w.Code != http.StatusNotFound {
			t.Errorf("expected %d, got %d", http.StatusNotFound, w.Code)
		}
	})
}
