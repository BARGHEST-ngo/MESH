package api_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/BARGHEST-ngo/MESH/provisioning/internal/api"
	"github.com/BARGHEST-ngo/MESH/provisioning/internal/state"
)

const testAPIKey = "test-key"

func newTestRouter(t *testing.T) http.Handler {
	t.Helper()
	reg, err := state.New(filepath.Join(t.TempDir(), "state.json"), 7001, 7010)
	if err != nil {
		t.Fatal(err)
	}
	return api.NewRouter(testAPIKey, reg)
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

func TestPostDeploymentPortExhaustion(t *testing.T) {
	reg, err := state.New(filepath.Join(t.TempDir(), "state.json"), 7001, 7001)
	if err != nil {
		t.Fatal(err)
	}
	router := api.NewRouter(testAPIKey, reg)

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
