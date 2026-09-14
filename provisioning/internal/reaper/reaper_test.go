package reaper_test

import (
	"context"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/BARGHEST-ngo/MESH/provisioning/internal/reaper"
	"github.com/BARGHEST-ngo/MESH/provisioning/internal/state"
)

const (
	minPort = 7001
	maxPort = 7010
)

type mockContainerService struct {
	mu        sync.Mutex
	failSlugs map[string]bool
	stopped   []string
}

func (m *mockContainerService) Start(state.Deployment, string) error { return nil }

func (m *mockContainerService) Stop(slug string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.stopped = append(m.stopped, slug)
	if m.failSlugs[slug] {
		return errSimulatedStopFailure
	}
	return nil
}

func (m *mockContainerService) stoppedSlugs() []string {
	m.mu.Lock()
	defer m.mu.Unlock()

	out := make([]string, 0)
	return append(out, m.stopped...)
}

var errSimulatedStopFailure = &simulatedError{"simulated stop failure"}

type simulatedError struct{ msg string }

func (e *simulatedError) Error() string { return e.msg }

func newTestRegistry(t *testing.T, defaultTTL time.Duration) *state.Registry {
	t.Helper()
	reg, err := state.New(filepath.Join(t.TempDir(), "state.json"), minPort, maxPort, defaultTTL)
	if err != nil {
		t.Fatal(err)
	}
	return reg
}

func runReaper(registry *state.Registry, service *mockContainerService, interval time.Duration) (stop func()) {
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() {
		reaper.Run(ctx, registry, service, interval)
		close(done)
	}()
	return func() {
		cancel()
		<-done
	}
}

func TestRemovesExpiredDeployment(t *testing.T) {
	r := newTestRegistry(t, time.Duration(-1)*time.Second)
	d, err := r.AllocatePort("test-slug", "test-owner", 1, nil)
	if err != nil {
		t.Fatal(err)
	}
	service := mockContainerService{
		failSlugs: make(map[string]bool),
		stopped:   make([]string, 0),
	}
	stopFunc := runReaper(r, &service, time.Second)
	defer stopFunc()

	deadline := time.Now().Add(3 * time.Second)
	pollInterval := 50 * time.Millisecond
	deploymentReaped := false
	for time.Now().Before(deadline) {
		if _, ok := r.Get(d.Slug); !ok {
			deploymentReaped = true
			break
		}
		time.Sleep(pollInterval)
	}

	if !deploymentReaped {
		t.Fatal("expected expired deployment to be removed within deadline")
	}

	stopped := service.stoppedSlugs()
	if len(stopped) != 1 || stopped[0] != d.Slug {
		t.Errorf("expected %s to have been stopped, got %v", d.Slug, stopped)
	}
}

func TestDoesNotReapUnexpiredDeployment(t *testing.T) {
	r := newTestRegistry(t, 30*time.Second)
	d, err := r.AllocatePort("test-slug", "test-owner", 1, nil)
	if err != nil {
		t.Fatal(err)
	}
	service := mockContainerService{
		failSlugs: make(map[string]bool),
		stopped:   make([]string, 0),
	}
	stopFunc := runReaper(r, &service, time.Second)
	defer stopFunc()

	deadline := time.Now().Add(3 * time.Second)
	pollInterval := 50 * time.Millisecond
	deploymentReaped := false
	for time.Now().Before(deadline) {
		if _, ok := r.Get(d.Slug); !ok {
			deploymentReaped = true
			break
		}
		time.Sleep(pollInterval)
	}

	if deploymentReaped {
		t.Error("expected deployment to not be removed")
	}
}

func TestStopsWhenContextCancelled(t *testing.T) {
	r := newTestRegistry(t, 30*time.Second)
	service := mockContainerService{
		failSlugs: make(map[string]bool),
		stopped:   make([]string, 0),
	}
	stopFunc := runReaper(r, &service, time.Minute)

	stopped := make(chan struct{})
	go func() {
		stopFunc()
		close(stopped)
	}()

	select {
	case <-stopped:
	case <-time.After(1 * time.Second):
		t.Fatal("stopFunc did not return within 1s - reaper should have stopped")
	}
}
