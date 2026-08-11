package state

// State Registry tracks the internal state of deployments
// Primarily tracks allocated ports for all started containers

import (
	"fmt"
	"sync"
	"time"
)

type Deployment struct {
	Slug      string    `json:"slug"`
	FrpsPort  int       `json:"frps_port"`
	CreatedAt time.Time `json:"created_at"`
	OwnerID   string    `json:"owner_id"`
}

type registryState struct {
	Deployments map[string]Deployment `json:"deployments"`
}

type Registry struct {
	mu      sync.Mutex
	state   registryState
	path    string
	portMin int
	portMax int
}

func New(path string, portMin, portMax int) (*Registry, error) {
	r := &Registry{
		path:    path,
		portMin: portMin,
		portMax: portMax,
		state:   registryState{Deployments: make(map[string]Deployment)},
	}
	if err := r.load(); err != nil {
		return nil, err
	}
	return r, nil
}

func (r *Registry) AllocatePort(slug, ownerID string, maxConcurrent int) (Deployment, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	used := make(map[int]bool)
	var count int
	for _, d := range r.state.Deployments {
		used[d.FrpsPort] = true
		if d.OwnerID == ownerID {
			count++
		}
	}

	if maxConcurrent != 0 && count >= maxConcurrent {
		return Deployment{}, fmt.Errorf("max deployments reached")
	}

	for port := r.portMin; port <= r.portMax; port++ {
		if !used[port] {
			d := Deployment{
				Slug:      slug,
				FrpsPort:  port,
				CreatedAt: time.Now().UTC(),
				OwnerID:   ownerID,
			}
			r.state.Deployments[slug] = d
			return d, r.save()
		}
	}

	return Deployment{}, fmt.Errorf("no available port")
}

func (r *Registry) Release(slug string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.state.Deployments[slug]; !ok {
		// TODO: Should this silently error?
		return fmt.Errorf("deployment %s not found", slug)
	}
	delete(r.state.Deployments, slug)
	return r.save()
}

func (r *Registry) Get(slug string) (Deployment, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()

	d, ok := r.state.Deployments[slug]
	return d, ok
}

func (r *Registry) load() error {
	return loadJSON(r.path, &r.state)
}

func (r *Registry) save() error {
	return saveJSON(r.path, r.state)
}
