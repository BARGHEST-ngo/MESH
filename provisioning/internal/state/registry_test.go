package state_test

import (
	"path/filepath"
	"testing"

	"github.com/BARGHEST-ngo/MESH/provisioning/internal/state"
)

const (
	testSlug = "test-slug"
	minPort  = 7000
	maxPort  = 7010
)

func defaultTestRegistry(t *testing.T) *state.Registry {
	t.Helper()
	reg, err := state.New(filepath.Join(t.TempDir(), "state.json"), minPort, maxPort)
	if err != nil {
		t.Fatal(err)
	}

	return reg
}

func TestAllocatePort(t *testing.T) {
	reg := defaultTestRegistry(t)

	created, err := reg.AllocatePort(testSlug, "test-owner", 0)
	if err != nil {
		t.Errorf("expected valid allocation")
	}

	if created.FrpsPort < minPort || created.FrpsPort > maxPort {
		t.Errorf("port allocated outside of permitted range: %d", created.FrpsPort)
	}

	found, ok := reg.Get(created.Slug)
	if !ok {
		t.Fatal("expected deployment to exist")
	}
	if found.FrpsPort != created.FrpsPort {
		t.Errorf("expected port %d, got %d", created.FrpsPort, found.FrpsPort)
	}
	if found.OwnerID != "test-owner" {
		t.Errorf("expected owner %q, got %q", "test-owner", found.OwnerID)
	}

	if _, ok := reg.Get("some-random-slug"); ok {
		t.Error("retrieved a slug that was never allocated")
	}
}

func TestRelease(t *testing.T) {
	reg := defaultTestRegistry(t)
	d, err := reg.AllocatePort(testSlug, "test-owner", 0)
	if err != nil {
		t.Fatal(err)
	}

	if err := reg.Release(d.Slug); err != nil {
		t.Fatal(err)
	}

	if _, ok := reg.Get(d.Slug); ok {
		t.Errorf("%s expected to be released", d.Slug)
	}

	if err := reg.Release(d.Slug); err == nil {
		t.Errorf("%s expected to already be released", d.Slug)
	}

	if err := reg.Release("random-slug"); err == nil {
		t.Error("released slug that was never allocated")
	}
}

func TestMaxConcurrent(t *testing.T) {
	t.Run("single-owner", func(t *testing.T) {
		const maxDeploys = 1
		reg := defaultTestRegistry(t)
		if _, err := reg.AllocatePort("testSlug-A", "test-owner", maxDeploys); err != nil {
			t.Fatal(err)
		}

		if _, err := reg.AllocatePort("testSlug-B", "test-owner", maxDeploys); err == nil {
			t.Errorf("deployed more than %d max deployments", maxDeploys)
		}
	})

	t.Run("multiple-owners", func(t *testing.T) {
		const maxDeploys = 1
		reg := defaultTestRegistry(t)
		if _, err := reg.AllocatePort("testSlug-A", "test-owner-0", maxDeploys); err != nil {
			t.Fatal(err)
		}

		if _, err := reg.AllocatePort("testSlug-B", "test-owner-1", maxDeploys); err != nil {
			t.Fatal(err)
		}

		if _, err := reg.AllocatePort("testSlug-C", "test-owner-0", maxDeploys); err == nil {
			t.Errorf("test-owner-0 deployed more than %d max deployments", maxDeploys)
		}

		if _, err := reg.AllocatePort("testSlug-D", "test-owner-1", maxDeploys); err == nil {
			t.Errorf("test-owner-1 deployed more than %d max deployments", maxDeploys)
		}
	})
}

func TestRegistryPersistence(t *testing.T) {
	path := filepath.Join(t.TempDir(), "state.json")
	reg, err := state.New(path, minPort, maxPort)
	if err != nil {
		t.Fatal(err)
	}
	created, err := reg.AllocatePort(testSlug, "test-owner", 0)
	if err != nil {
		t.Fatal(err)
	}

	// create a new registry with the same file path as before to simulate a restart
	reg2, err := state.New(path, minPort, maxPort)
	if err != nil {
		t.Fatal(err)
	}
	found, ok := reg2.Get(testSlug)
	if !ok {
		t.Errorf("port allocation did not persist between restarts")
	}
	if found.FrpsPort != created.FrpsPort {
		t.Errorf("expected port %d, got %d", created.FrpsPort, found.FrpsPort)
	}
	if found.OwnerID != "test-owner" {
		t.Errorf("expected owner %q, got %q", "test-owner", found.OwnerID)
	}
}
