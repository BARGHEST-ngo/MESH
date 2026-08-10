package workerclient

import (
	"fmt"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/BARGHEST-ngo/MESH/provisioning/internal/state"
	"github.com/BARGHEST-ngo/MESH/provisioning/internal/workerapi"
)

const testToken = "test-worker-token"

type mockRunner struct {
	startErr    error
	stopErr     error
	startedWith state.Deployment
	stoppedSlug string
}

func (m *mockRunner) Start(d state.Deployment) error {
	m.startedWith = d
	return m.startErr
}

func (m *mockRunner) Stop(slug string) error {
	m.stoppedSlug = slug
	return m.stopErr
}

func newTestServer(runner *mockRunner) *httptest.Server {
	return httptest.NewServer(workerapi.NewWorkerRouter(runner, testToken))
}

func TestClientStart(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		runner := &mockRunner{}
		srv := newTestServer(runner)
		defer srv.Close()

		client := New(srv.URL, testToken)
		d := state.Deployment{Slug: "abc123def4", Token: "tok", FrpsPort: 7001, CreatedAt: time.Now().UTC(), OwnerID: "owner-a"}
		if err := client.Start(d); err != nil {
			t.Fatalf("expected success, got %v", err)
		}

		if runner.startedWith.Slug != d.Slug || runner.startedWith.OwnerID != d.OwnerID {
			t.Errorf("worker received unexpected deployment: %+v", runner.startedWith)
		}
	})

	t.Run("worker error propagates", func(t *testing.T) {
		runner := &mockRunner{startErr: fmt.Errorf("docker exploded")}
		srv := newTestServer(runner)
		defer srv.Close()

		client := New(srv.URL, testToken)
		err := client.Start(state.Deployment{Slug: "abc123def4"})
		if err == nil {
			t.Fatal("expected error")
		}
		if !strings.Contains(err.Error(), "500") {
			t.Errorf("expected error to mention 500 status, got %q", err.Error())
		}
	})

	t.Run("wrong token", func(t *testing.T) {
		runner := &mockRunner{}
		srv := newTestServer(runner)
		defer srv.Close()

		client := New(srv.URL, "wrong-token")
		err := client.Start(state.Deployment{Slug: "abc123def4"})
		if err == nil {
			t.Fatal("expected error")
		}
		if !strings.Contains(err.Error(), "401") {
			t.Errorf("expected error to mention 401 status, got %q", err.Error())
		}
	})
}

func TestClientStop(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		runner := &mockRunner{}
		srv := newTestServer(runner)
		defer srv.Close()

		client := New(srv.URL, testToken)
		if err := client.Stop("abc123def4"); err != nil {
			t.Fatalf("expected success, got %v", err)
		}
		if runner.stoppedSlug != "abc123def4" {
			t.Errorf("expected worker to receive slug abc123def4, got %q", runner.stoppedSlug)
		}
	})

	t.Run("worker error propagates", func(t *testing.T) {
		runner := &mockRunner{stopErr: fmt.Errorf("docker exploded")}
		srv := newTestServer(runner)
		defer srv.Close()

		client := New(srv.URL, testToken)
		err := client.Stop("abc123def4")
		if err == nil {
			t.Fatal("expected error")
		}
		if !strings.Contains(err.Error(), "500") {
			t.Errorf("expected error to mention 500 status, got %q", err.Error())
		}
	})
}
