package workerapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/BARGHEST-ngo/MESH/provisioning/internal/state"
)

const testWorkerToken = "test-worker-token"

type mockRunner struct {
	startErr error
	stopErr  error

	startedWith state.Deployment
	stoppedSlug string
	startCalled bool
	stopCalled  bool
}

func (m *mockRunner) Start(d state.Deployment) error {
	m.startCalled = true
	m.startedWith = d
	return m.startErr
}

func (m *mockRunner) Stop(slug string) error {
	m.stopCalled = true
	m.stoppedSlug = slug
	return m.stopErr
}

func newAuthedRequest(method, target string, body []byte) *http.Request {
	req := httptest.NewRequest(method, target, bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+testWorkerToken)
	return req
}

func TestHandleStart(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		runner := &mockRunner{}
		router := NewWorkerRouter(runner, testWorkerToken)

		d := state.Deployment{Slug: "abc123def4", Token: "tok", FrpsPort: 7001, CreatedAt: time.Now().UTC(), OwnerID: "owner-a"}
		b, err := json.Marshal(d)
		if err != nil {
			t.Fatal(err)
		}

		req := newAuthedRequest(http.MethodPost, "/containers", b)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusNoContent {
			t.Errorf("expected %d, got %d", http.StatusNoContent, w.Code)
		}
		if !runner.startCalled {
			t.Fatal("expected runner.Start to be called")
		}
		if runner.startedWith.Slug != d.Slug || runner.startedWith.OwnerID != d.OwnerID {
			t.Errorf("runner received unexpected deployment: %+v", runner.startedWith)
		}
	})

	t.Run("missing token", func(t *testing.T) {
		runner := &mockRunner{}
		router := NewWorkerRouter(runner, testWorkerToken)

		req := httptest.NewRequest(http.MethodPost, "/containers", bytes.NewReader([]byte(`{}`)))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("expected %d, got %d", http.StatusUnauthorized, w.Code)
		}
		if runner.startCalled {
			t.Error("runner should not be called on auth failure")
		}
	})

	t.Run("wrong token", func(t *testing.T) {
		runner := &mockRunner{}
		router := NewWorkerRouter(runner, testWorkerToken)

		req := httptest.NewRequest(http.MethodPost, "/containers", bytes.NewReader([]byte(`{}`)))
		req.Header.Set("Authorization", "Bearer wrong-token")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("expected %d, got %d", http.StatusUnauthorized, w.Code)
		}
	})

	t.Run("malformed json", func(t *testing.T) {
		runner := &mockRunner{}
		router := NewWorkerRouter(runner, testWorkerToken)

		req := newAuthedRequest(http.MethodPost, "/containers", []byte("not json"))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected %d, got %d", http.StatusBadRequest, w.Code)
		}
		if runner.startCalled {
			t.Error("runner should not be called on bad request body")
		}
	})

	t.Run("runner error", func(t *testing.T) {
		runner := &mockRunner{startErr: fmt.Errorf("docker exploded")}
		router := NewWorkerRouter(runner, testWorkerToken)

		b, err := json.Marshal(state.Deployment{Slug: "abc123def4"})
		if err != nil {
			t.Fatal(err)
		}
		req := newAuthedRequest(http.MethodPost, "/containers", b)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("expected %d, got %d", http.StatusInternalServerError, w.Code)
		}
	})
}

func TestHandleStop(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		runner := &mockRunner{}
		router := NewWorkerRouter(runner, testWorkerToken)

		req := newAuthedRequest(http.MethodDelete, "/containers/abc123def4", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusNoContent {
			t.Errorf("expected %d, got %d", http.StatusNoContent, w.Code)
		}
		if !runner.stopCalled || runner.stoppedSlug != "abc123def4" {
			t.Errorf("expected runner.Stop called with abc123def4, got called=%v slug=%q", runner.stopCalled, runner.stoppedSlug)
		}
	})

	t.Run("invalid slug traversal", func(t *testing.T) {
		runner := &mockRunner{}
		router := NewWorkerRouter(runner, testWorkerToken)

		req := newAuthedRequest(http.MethodDelete, "/containers/%2e%2e", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected %d, got %d", http.StatusBadRequest, w.Code)
		}
		if runner.stopCalled {
			t.Error("runner should not be called for an invalid slug")
		}
	})

	t.Run("wrong token", func(t *testing.T) {
		runner := &mockRunner{}
		router := NewWorkerRouter(runner, testWorkerToken)

		req := httptest.NewRequest(http.MethodDelete, "/containers/abc123def4", nil)
		req.Header.Set("Authorization", "Bearer wrong-token")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("expected %d, got %d", http.StatusUnauthorized, w.Code)
		}
	})

	t.Run("runner error", func(t *testing.T) {
		runner := &mockRunner{stopErr: fmt.Errorf("docker exploded")}
		router := NewWorkerRouter(runner, testWorkerToken)

		req := newAuthedRequest(http.MethodDelete, "/containers/abc123def4", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("expected %d, got %d", http.StatusInternalServerError, w.Code)
		}
	})
}
