package api_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/BARGHEST-ngo/MESH/provisioning/internal/api"
)

func TestAuthRequest(t *testing.T) {
	dummy := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	handler := api.AuthRequest(newTestKeyStore(t), dummy)

	cases := []struct {
		name     string
		token    string
		expected int
	}{
		{"no token", "", http.StatusUnauthorized},
		{"wrong token", "wrong-key", http.StatusUnauthorized},
		{"valid token", "test-key", http.StatusOK},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/dummy", nil)
			if tc.token != "" {
				req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tc.token))
			}
			w := httptest.NewRecorder()
			handler.ServeHTTP(w, req)

			if w.Code != tc.expected {
				t.Errorf("expected %d, got %d", tc.expected, w.Code)
			}
		})
	}
}
