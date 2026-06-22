package state

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"
)

type Deployment struct {
	Slug      string    `json:"slug"`
	Token     string    `json:"token"`
	FrpsPort  int       `json:"frps_port"`
	CreatedAt time.Time `json:"created_at"`
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

func (r *Registry) AllocatePort(slug, token string) (int, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	used := make(map[int]bool)
	for _, d := range r.state.Deployments {
		used[d.FrpsPort] = true
	}

	for port := r.portMin; port <= r.portMax; port++ {
		if !used[port] {
			r.state.Deployments[slug] = Deployment{
				Slug:      slug,
				Token:     token,
				FrpsPort:  port,
				CreatedAt: time.Now().UTC(),
			}
			return port, r.save()
		}
	}

	return 0, fmt.Errorf("no available port")
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
	data, err := os.ReadFile(r.path)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("error reading depolyments file: %w", err)
	}
	return json.Unmarshal(data, &r.state)
}

func (r *Registry) save() error {
	data, err := json.MarshalIndent(r.state, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal error: %w", err)
	}

	tmp := r.path + ".tmp"
	if err := os.WriteFile(tmp, data, 0600); err != nil {
		return fmt.Errorf("write file error: %w", err)
	}
	return os.Rename(tmp, r.path)
}
