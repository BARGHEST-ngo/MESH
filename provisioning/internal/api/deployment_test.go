package api_test

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/BARGHEST-ngo/MESH/provisioning/internal/api"
	"github.com/BARGHEST-ngo/MESH/provisioning/internal/state"
)

const (
	testAPIKey = "test-key"
	minPort    = 7001
	maxPort    = 7010
)

type mockContainerService struct{}

func (mockContainerService) Start(state.Deployment) error { return nil }
func (mockContainerService) Stop(string) error            { return nil }

type failOnceMock struct{ calls int }

func (m *failOnceMock) Start(state.Deployment) error {
	m.calls++
	if m.calls == 1 {
		return fmt.Errorf("simulated start failure")
	}
	return nil
}
func (m *failOnceMock) Stop(string) error { return nil }

func defaultTestKey() state.APIKey {
	hash := sha256.Sum256([]byte(testAPIKey))
	created := time.Now()
	expires := created.Add(time.Hour)
	return state.APIKey{
		ID:            "test-owner",
		Label:         "test",
		HashHex:       hex.EncodeToString(hash[:]),
		MaxConcurrent: 0,
		Revoked:       false,
		CreatedAt:     created,
		ExpiresAt:     &expires,
	}
}

func newTestRouterWIthKey(t *testing.T, key state.APIKey) http.Handler {
	t.Helper()
	reg, err := state.New(filepath.Join(t.TempDir(), "state.json"), 7001, 7010)
	if err != nil {
		t.Fatal(err)
	}
	return api.NewRouter(newTestKeyStoreWithKey(t, key), reg, "some_frps_image", api.WithContainerService(mockContainerService{}))
}

func newTestRouter(t *testing.T) http.Handler {
	t.Helper()
	reg, err := state.New(filepath.Join(t.TempDir(), "state.json"), 7001, 7010)
	if err != nil {
		t.Fatal(err)
	}
	return api.NewRouter(newTestKeyStore(t), reg, "some_frps_image", api.WithContainerService(mockContainerService{}))
}

func newTestKeyStoreWithKey(t *testing.T, key state.APIKey) *state.KeyStore {
	t.Helper()
	data, err := json.Marshal(struct {
		Keys []state.APIKey `json:"keys"`
	}{Keys: []state.APIKey{key}})
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(t.TempDir(), "keys.json")
	if err := os.WriteFile(path, data, 0600); err != nil {
		t.Fatal(err)
	}
	ks, err := state.NewKeyStore(path)
	if err != nil {
		t.Fatal(err)
	}
	return ks
}

func newTestKeyStore(t *testing.T) *state.KeyStore {
	t.Helper()
	return newTestKeyStoreWithKey(t, defaultTestKey())
}

func TestHealth(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()
	newTestRouter(t).ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("cannot reach health endpoint")
	}
}

func TestPostDeployment(t *testing.T) {
	cases := []struct {
		name     string
		token    string
		expected int
		method   string
	}{
		{"no token", "", http.StatusUnauthorized, http.MethodPost},
		{"wrong token", "wrong-token", http.StatusUnauthorized, http.MethodPost},
		{"valid token", testAPIKey, http.StatusCreated, http.MethodPost},
		{"invalid method", testAPIKey, http.StatusMethodNotAllowed, http.MethodGet},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(tc.method, "/deployment", nil)
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tc.token))
			w := httptest.NewRecorder()
			newTestRouter(t).ServeHTTP(w, req)

			if w.Code != tc.expected {
				t.Errorf("expected %d, got %d", tc.expected, w.Code)
			}
		})
	}
}

func TestPostDeploymentStructure(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/deployment", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", testAPIKey))
	w := httptest.NewRecorder()
	newTestRouter(t).ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("failed to create deployment")
	}

	var d api.DeploymentResponse
	if err := json.NewDecoder(w.Body).Decode(&d); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(d.Slug) != 10 {
		t.Errorf("expected valid slug")
	}

	if d.FrpsPort < minPort || d.FrpsPort > maxPort {
		t.Errorf("expected valid port")
	}
}

func TestPostDeploymentPortExhaustion(t *testing.T) {
	reg, err := state.New(filepath.Join(t.TempDir(), "state.json"), 7001, 7001)
	if err != nil {
		t.Fatal(err)
	}
	router := api.NewRouter(newTestKeyStore(t), reg, "some_frps_image", api.WithContainerService(mockContainerService{}))

	makeRequest := func() int {
		req := httptest.NewRequest(http.MethodPost, "/deployment", nil)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", testAPIKey))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		return w.Code
	}

	if code := makeRequest(); code != http.StatusCreated {
		t.Fatalf("first request expected %d, got %d", http.StatusCreated, code)
	}

	if code := makeRequest(); code != http.StatusInternalServerError {
		t.Errorf("second request expected %d, got %d", http.StatusInternalServerError, code)
	}
}

func TestStartFailureRollsBackPort(t *testing.T) {
	reg, err := state.New(filepath.Join(t.TempDir(), "state.json"), 7001, 7001)
	if err != nil {
		t.Fatal(err)
	}
	mock := &failOnceMock{}
	router := api.NewRouter(newTestKeyStore(t), reg, "some_frps_image", api.WithContainerService(mock))

	post := func() int {
		req := httptest.NewRequest(http.MethodPost, "/deployment", nil)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", testAPIKey))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		return w.Code
	}

	if code := post(); code != http.StatusInternalServerError {
		t.Fatalf("expected start failure to return 500, got %d", code)
	}
	if code := post(); code != http.StatusCreated {
		t.Errorf("expected port to be freed after rollback, got %d", code)
	}
}

func TestConcurrentDeployments(t *testing.T) {
	var mu sync.Mutex
	var deployments []api.DeploymentResponse
	var codes []int

	reg, err := state.New(filepath.Join(t.TempDir(), "state.json"), 7001, 7010)
	if err != nil {
		t.Fatalf("failed to create registry")
	}
	router := api.NewRouter(newTestKeyStore(t), reg, "some_frps_image", api.WithContainerService(mockContainerService{}))

	var wg sync.WaitGroup
	for range 10 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			req := httptest.NewRequest(http.MethodPost, "/deployment", nil)
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", testAPIKey))
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			mu.Lock()
			codes = append(codes, w.Code)
			if w.Code == http.StatusCreated {
				var created api.DeploymentResponse
				json.NewDecoder(w.Body).Decode(&created)
				deployments = append(deployments, created)
			}
			mu.Unlock()
		}()
	}
	wg.Wait()

	for _, code := range codes {
		if code != http.StatusCreated {
			t.Fatalf("expected %d, got %d", http.StatusCreated, code)
		}
	}

	cases := []struct {
		name  string
		check func([]api.DeploymentResponse) error
	}{
		{
			"duplicate-ports",
			func(results []api.DeploymentResponse) error {
				seen := make(map[int]struct{})
				for _, d := range results {
					if _, exists := seen[d.FrpsPort]; exists {
						return fmt.Errorf("duplicate port allocated: %d", d.FrpsPort)
					}
					seen[d.FrpsPort] = struct{}{}
				}
				return nil
			},
		},
		{
			"duplicate-slugs",
			func(results []api.DeploymentResponse) error {
				seen := make(map[string]struct{})
				for _, d := range results {
					if _, exists := seen[d.Slug]; exists {
						return fmt.Errorf("duplicate slug allocated: %s", d.Slug)
					}
					seen[d.Slug] = struct{}{}
				}
				return nil
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if err := tc.check(deployments); err != nil {
				t.Errorf("failed: %v", err)
			}
		})
	}
}

func TestApiKeyStates(t *testing.T) {
	//keyHash := sha256.Sum256([]byte(testAPIKey))

	post := func(router http.Handler) int {
		req := httptest.NewRequest(http.MethodPost, "/deployment", nil)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", testAPIKey))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		return w.Code
	}

	t.Run("expired-key", func(t *testing.T) {
		expires := time.Now().Add(-time.Hour)
		key := defaultTestKey()
		key.ExpiresAt = &expires

		if code := post(newTestRouterWIthKey(t, key)); code != http.StatusUnauthorized {
			t.Errorf("expected %d, got %d", http.StatusUnauthorized, code)
		}
	})

	t.Run("revoked-key", func(t *testing.T) {
		key := defaultTestKey()
		key.Revoked = true

		if code := post(newTestRouterWIthKey(t, key)); code != http.StatusUnauthorized {
			t.Errorf("expected %d, got %d", http.StatusUnauthorized, code)
		}
	})

	t.Run("over-max-concurrent-deploys", func(t *testing.T) {
		key := defaultTestKey()
		key.MaxConcurrent = 2
		reg, err := state.New(filepath.Join(t.TempDir(), "state.json"), 7001, 7010)
		if err != nil {
			t.Fatal(err)
		}
		router := api.NewRouter(newTestKeyStoreWithKey(t, key), reg, "some-frps-image", api.WithContainerService(mockContainerService{}))
		// 2 valid deployments, fail on the 3rd
		if code := post(router); code != http.StatusCreated {
			t.Errorf("expected %d, got %d", http.StatusCreated, code)
		}
		if code := post(router); code != http.StatusCreated {
			t.Errorf("expected %d, got %d", http.StatusCreated, code)
		}
		if code := post(router); code != http.StatusInternalServerError {
			t.Errorf("expected %d, got %d", http.StatusInternalServerError, code)
		}
	})
}

func TestDeleteDeployment(t *testing.T) {
	router := newTestRouter(t)
	req := httptest.NewRequest(http.MethodPost, "/deployment", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", testAPIKey))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatal("failed to create deployment")
	}
	var created api.DeploymentResponse
	json.NewDecoder(w.Body).Decode(&created)

	cases := []struct {
		name     string
		slug     string
		expected int
	}{
		{"ok", created.Slug, http.StatusOK},
		{"slug does not exist", "wrong-slug", http.StatusNotFound},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/deployment/%s", tc.slug), nil)
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", testAPIKey))
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			if w.Code != tc.expected {
				t.Errorf("expected %d, got %d", tc.expected, w.Code)
			}
		})
	}
}
