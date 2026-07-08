package api_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"sync"
	"testing"

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

func newTestRouter(t *testing.T) http.Handler {
	t.Helper()
	reg, err := state.New(filepath.Join(t.TempDir(), "state.json"), 7001, 7010)
	if err != nil {
		t.Fatal(err)
	}
	return api.NewRouter(testAPIKey, reg, "some_frps_image", api.WithContainerService(mockContainerService{}))
}

func validTestDeployment(t *testing.T) api.DeploymentResponse {
	req := httptest.NewRequest(http.MethodPost, "/deployment", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", testAPIKey))
	w := httptest.NewRecorder()
	newTestRouter(t).ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("failed to create deployment")
	}
	var resp api.DeploymentResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	return resp
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
	cases := []struct {
		name string
		test func(api.DeploymentResponse) bool
	}{
		{"valid slug", func(d api.DeploymentResponse) bool {
			return len(d.Slug) == 10
		}},
		{"valid port", func(d api.DeploymentResponse) bool {
			return d.FrpsPort >= minPort && d.FrpsPort <= maxPort
		}},
	}

	d := validTestDeployment(t)
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if ok := tc.test(d); !ok {
				t.Errorf("%s failed", tc.name)
			}
		})
	}
}

func TestPostDeploymentPortExhaustion(t *testing.T) {
	reg, err := state.New(filepath.Join(t.TempDir(), "state.json"), 7001, 7001)
	if err != nil {
		t.Fatal(err)
	}
	router := api.NewRouter(testAPIKey, reg, "some_frps_image", api.WithContainerService(mockContainerService{}))

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
	router := api.NewRouter(testAPIKey, reg, "some_frps_image", api.WithContainerService(mock))

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
	router := api.NewRouter(testAPIKey, reg, "some_frps_image", api.WithContainerService(mockContainerService{}))

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
