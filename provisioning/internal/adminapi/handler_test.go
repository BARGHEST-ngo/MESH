package adminapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"
	"time"

	"github.com/BARGHEST-ngo/MESH/provisioning/internal/state"
)

const testAdminToken = "test-admin-token"

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
	return NewAdminRouter(newTestKeyStore(t), testAdminToken)
}

func newAuthedRequest(method, target string, body io.Reader) *http.Request {
	req := httptest.NewRequest(method, target, body)
	req.Header.Set("Authorization", "Bearer "+testAdminToken)
	return req
}

func createPatchTestKey(t *testing.T, router http.Handler) createKeyResponse {
	t.Helper()
	requestData := createKeyRequest{
		OwnerID:       "test-owner",
		Label:         "test-label",
		MaxConcurrent: 1,
		TTLHours:      1,
	}
	b, err := json.Marshal(&requestData)
	if err != nil {
		t.Fatal(err)
	}

	createReq := newAuthedRequest(http.MethodPost, "/keys", bytes.NewBuffer(b))
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

func listKeys(t *testing.T, router http.Handler) getKeysResponse {
	t.Helper()
	getReq := newAuthedRequest(http.MethodGet, "/keys", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, getReq)
	if w.Code != http.StatusOK {
		t.Fatal("expected success")
	}
	var response getKeysResponse
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatal(err)
	}

	return response
}

func TestCreateKey(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		const ownerId = "test-owner"
		const label = "test-label"
		requestData := createKeyRequest{
			OwnerID:       ownerId,
			Label:         label,
			MaxConcurrent: 1,
			TTLHours:      0,
		}
		b, err := json.Marshal(&requestData)
		if err != nil {
			t.Fatal(err)
		}

		req := newAuthedRequest(http.MethodPost, "/keys", bytes.NewBuffer(b))
		w := httptest.NewRecorder()
		router := newTestAdminRouter(t)
		router.ServeHTTP(w, req)
		if w.Code != http.StatusCreated {
			t.Errorf("expected %d, got %d", http.StatusCreated, w.Code)
		}

		var created createKeyResponse
		if err := json.Unmarshal(w.Body.Bytes(), &created); err != nil {
			t.Fatal(err)
		}

		if created.ID == "" || created.Key == "" || created.Label != label || created.OwnerID != ownerId || created.ExpiresAt != nil {
			t.Error("invalid response")
		}

		listKeysResponse := listKeys(t, router)
		if _, ok := listKeysResponse.Keys[created.ID]; !ok {
			t.Errorf("created key is not present in GET keys")
		}
	})

	t.Run("valid-with-ttl", func(t *testing.T) {
		const ownerId = "test-owner"
		const label = "test-label"
		requestData := createKeyRequest{
			OwnerID:       ownerId,
			Label:         label,
			MaxConcurrent: 1,
			TTLHours:      1,
		}
		b, err := json.Marshal(&requestData)
		if err != nil {
			t.Fatal(err)
		}

		req := newAuthedRequest(http.MethodPost, "/keys", bytes.NewBuffer(b))
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

	t.Run("valid-with-deployment-ttl", func(t *testing.T) {
		const ownerId = "test-owner"
		const label = "test-label"
		const deploymentTTLHours = 6

		requestData := createKeyRequest{
			OwnerID:            ownerId,
			Label:              label,
			MaxConcurrent:      1,
			DeploymentTTLHours: deploymentTTLHours,
		}
		b, err := json.Marshal(&requestData)
		if err != nil {
			t.Fatal(err)
		}

		req := newAuthedRequest(http.MethodPost, "/keys", bytes.NewBuffer(b))
		w := httptest.NewRecorder()
		newTestAdminRouter(t).ServeHTTP(w, req)
		if w.Code != http.StatusCreated {
			t.Errorf("expected %d, got %d", http.StatusCreated, w.Code)
		}

		var response createKeyResponse
		if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
			t.Fatal(err)
		}

		if response.DeploymentTTLHours == nil {
			t.Fatal("expected DeploymentTTL to be set")
		}
		if *response.DeploymentTTLHours != deploymentTTLHours {
			t.Errorf("expected DeploymentTTL %v, got %v", deploymentTTLHours, *response.DeploymentTTLHours)
		}
	})

	t.Run("valid-without-deployment-ttl", func(t *testing.T) {
		requestData := createKeyRequest{
			OwnerID:       "test-owner",
			Label:         "test-label",
			MaxConcurrent: 1,
		}
		b, err := json.Marshal(&requestData)
		if err != nil {
			t.Fatal(err)
		}

		req := newAuthedRequest(http.MethodPost, "/keys", bytes.NewBuffer(b))
		w := httptest.NewRecorder()
		newTestAdminRouter(t).ServeHTTP(w, req)
		if w.Code != http.StatusCreated {
			t.Errorf("expected %d, got %d", http.StatusCreated, w.Code)
		}
		var response createKeyResponse
		if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
			t.Fatal(err)
		}

		if response.DeploymentTTLHours != nil {
			t.Errorf("expected DeploymentTTL to be nil (use default), got %v", *response.DeploymentTTLHours)
		}
	})

	t.Run("empty-data", func(t *testing.T) {
		b, err := json.Marshal(&createKeyRequest{})
		if err != nil {
			t.Fatal(err)
		}

		req := newAuthedRequest(http.MethodPost, "/keys", bytes.NewBuffer(b))
		w := httptest.NewRecorder()
		newTestAdminRouter(t).ServeHTTP(w, req)
		if w.Code != http.StatusBadRequest {
			t.Errorf("expected %d, got %d", http.StatusBadRequest, w.Code)
		}
	})

	t.Run("malformed-json", func(t *testing.T) {
		b := []byte("not json")
		req := newAuthedRequest(http.MethodPost, "/keys", bytes.NewBuffer(b))
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

		req := newAuthedRequest(http.MethodPost, "/keys", bytes.NewBuffer(b))
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

		req := newAuthedRequest(http.MethodPost, "/keys", bytes.NewBuffer(b))
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

		req := newAuthedRequest(http.MethodPost, "/keys", bytes.NewBuffer(b))
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
			MaxConcurrent: 1,
			TTLHours:      0,
		}
		b, err := json.Marshal(&requestData)
		if err != nil {
			t.Fatal(err)
		}

		createReq := newAuthedRequest(http.MethodPost, "/keys", bytes.NewBuffer(b))
		w := httptest.NewRecorder()
		router := newTestAdminRouter(t)
		router.ServeHTTP(w, createReq)
		if w.Code != http.StatusCreated {
			t.Errorf("expected %d, got %d", http.StatusCreated, w.Code)
		}

		var created createKeyResponse
		if err := json.Unmarshal(w.Body.Bytes(), &created); err != nil {
			t.Fatal(err)
		}

		deleteReq := newAuthedRequest(http.MethodDelete, fmt.Sprintf("/keys/%s", created.ID), nil)
		w = httptest.NewRecorder()
		router.ServeHTTP(w, deleteReq)
		if w.Code != http.StatusOK {
			t.Fatalf("expected %d, got %d", http.StatusOK, w.Code)
		}

		listKeysResponse := listKeys(t, router)
		if k, ok := listKeysResponse.Keys[created.ID]; !ok || !k.Revoked {
			t.Errorf("created key is not revoked")
		}
	})

	t.Run("unknown-key", func(t *testing.T) {
		deleteReq := newAuthedRequest(http.MethodDelete, "/keys/unknown-id", nil)
		w := httptest.NewRecorder()
		newTestAdminRouter(t).ServeHTTP(w, deleteReq)
		if w.Code != http.StatusNotFound {
			t.Errorf("expected %d, got %d", http.StatusNotFound, w.Code)
		}
	})
}

func TestPatchKey(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		router := newTestAdminRouter(t)
		created := createPatchTestKey(t, router)
		newLabel := "new-label"
		maxConcurrent := 1
		patchData := updateKeyRequest{
			Label:         &newLabel,
			MaxConcurrent: &maxConcurrent,
		}
		b, err := json.Marshal(&patchData)
		if err != nil {
			t.Fatal(err)
		}
		patchReq := newAuthedRequest(http.MethodPatch, fmt.Sprintf("/keys/%s", created.ID), bytes.NewBuffer(b))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, patchReq)
		if w.Code != http.StatusOK {
			t.Fatalf("expected %d, got %d", http.StatusOK, w.Code)
		}

		listKeysResponse := listKeys(t, router)
		if k, ok := listKeysResponse.Keys[created.ID]; !ok {
			t.Fatalf("patched key is not present in GET keys")
		} else if k.Label != newLabel || k.MaxConcurrent != maxConcurrent || !k.ExpiresAt.Equal(*created.ExpiresAt) {
			t.Errorf("key has unexpected changes")
		}
	})

	t.Run("valid-no-change", func(t *testing.T) {
		router := newTestAdminRouter(t)
		created := createPatchTestKey(t, router)

		patchData := updateKeyRequest{}
		b, err := json.Marshal(&patchData)
		if err != nil {
			t.Fatal(err)
		}
		patchReq := newAuthedRequest(http.MethodPatch, fmt.Sprintf("/keys/%s", created.ID), bytes.NewBuffer(b))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, patchReq)
		if w.Code != http.StatusOK {
			t.Errorf("expected %d, got %d", http.StatusOK, w.Code)
		}

		listKeysResponse := listKeys(t, router)
		if k, ok := listKeysResponse.Keys[created.ID]; !ok {
			t.Fatalf("patched key is not present in GET keys")
		} else if k.Label != created.Label || k.MaxConcurrent != 1 || !k.ExpiresAt.Equal(*created.ExpiresAt) {
			t.Errorf("key has unexpected changes")
		}
	})

	t.Run("clear-expiry", func(t *testing.T) {
		router := newTestAdminRouter(t)
		created := createPatchTestKey(t, router)
		patchData := updateKeyRequest{
			ClearExpiry: true,
		}
		b, err := json.Marshal(&patchData)
		if err != nil {
			t.Fatal(err)
		}
		patchReq := newAuthedRequest(http.MethodPatch, fmt.Sprintf("/keys/%s", created.ID), bytes.NewBuffer(b))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, patchReq)
		if w.Code != http.StatusOK {
			t.Errorf("expected %d, got %d", http.StatusOK, w.Code)
		}

		listKeysResponse := listKeys(t, router)
		if k, ok := listKeysResponse.Keys[created.ID]; !ok {
			t.Fatalf("patched key is not present in GET keys")
		} else if k.Label != created.Label || k.MaxConcurrent != 1 || k.ExpiresAt != nil {
			t.Errorf("key has unexpected changes")
		}
	})

	t.Run("set-expiry", func(t *testing.T) {
		router := newTestAdminRouter(t)
		created := createPatchTestKey(t, router)
		newExpiry := time.Now().Add(time.Hour)
		patchData := updateKeyRequest{
			ExpiresAt: &newExpiry,
		}
		b, err := json.Marshal(&patchData)
		if err != nil {
			t.Fatal(err)
		}
		patchReq := newAuthedRequest(http.MethodPatch, fmt.Sprintf("/keys/%s", created.ID), bytes.NewBuffer(b))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, patchReq)
		if w.Code != http.StatusOK {
			t.Errorf("expected %d, got %d", http.StatusOK, w.Code)
		}

		listKeysResponse := listKeys(t, router)
		if k, ok := listKeysResponse.Keys[created.ID]; !ok {
			t.Fatalf("patched key is not present in GET keys")
		} else if k.Label != created.Label || k.MaxConcurrent != 1 || !k.ExpiresAt.Equal(newExpiry) {
			t.Errorf("key has unexpected changes")
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
		patchReq := newAuthedRequest(http.MethodPatch, fmt.Sprintf("/keys/%s", response.ID), bytes.NewBuffer(b))
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
		patchReq := newAuthedRequest(http.MethodPatch, fmt.Sprintf("/keys/%s", response.ID), bytes.NewBuffer(b))
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
		patchReq := newAuthedRequest(http.MethodPatch, "/keys/unknown-id", bytes.NewBuffer(b))
		w := httptest.NewRecorder()
		newTestAdminRouter(t).ServeHTTP(w, patchReq)
		if w.Code != http.StatusNotFound {
			t.Errorf("expected %d, got %d", http.StatusNotFound, w.Code)
		}
	})
}

func TestBodySizeLimit(t *testing.T) {
	t.Run("oversized body rejected", func(t *testing.T) {
		requestData := createKeyRequest{
			OwnerID:       "owner",
			Label:         string(bytes.Repeat([]byte("a"), maxBodyBytes)),
			MaxConcurrent: 1,
		}
		b, err := json.Marshal(&requestData)
		if err != nil {
			t.Fatal(err)
		}

		req := newAuthedRequest(http.MethodPost, "/keys", bytes.NewBuffer(b))
		w := httptest.NewRecorder()
		newTestAdminRouter(t).ServeHTTP(w, req)
		if w.Code != http.StatusBadRequest {
			t.Errorf("expected %d, got %d", http.StatusBadRequest, w.Code)
		}
	})

	t.Run("body at limit accepted", func(t *testing.T) {
		empty, err := json.Marshal(&createKeyRequest{OwnerID: "owner", MaxConcurrent: 1})
		if err != nil {
			t.Fatal(err)
		}
		overhead := len(empty)

		requestData := createKeyRequest{
			OwnerID:       "owner",
			Label:         string(bytes.Repeat([]byte("a"), maxBodyBytes-overhead)),
			MaxConcurrent: 1,
		}
		b, err := json.Marshal(&requestData)
		if err != nil {
			t.Fatal(err)
		}
		if len(b) != maxBodyBytes {
			t.Fatalf("test body is %d bytes, want exactly maxBodyBytes (%d)", len(b), maxBodyBytes)
		}

		req := newAuthedRequest(http.MethodPost, "/keys", bytes.NewBuffer(b))
		w := httptest.NewRecorder()
		newTestAdminRouter(t).ServeHTTP(w, req)
		if w.Code != http.StatusCreated {
			t.Errorf("expected %d, got %d", http.StatusCreated, w.Code)
		}
	})
}

func TestAdminAuth(t *testing.T) {
	router := newTestAdminRouter(t)

	cases := []struct {
		name   string
		header string
	}{
		{"no header", ""},
		{"malformed header", testAdminToken},
		{"wrong token", "Bearer wrong-token"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/keys", nil)
			if tc.header != "" {
				req.Header.Set("Authorization", tc.header)
			}
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			if w.Code != http.StatusUnauthorized {
				t.Errorf("expected %d, got %d", http.StatusUnauthorized, w.Code)
			}
		})
	}

	t.Run("valid token", func(t *testing.T) {
		req := newAuthedRequest(http.MethodGet, "/keys", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Errorf("expected %d, got %d", http.StatusOK, w.Code)
		}
	})
}
