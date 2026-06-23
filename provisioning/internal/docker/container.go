package docker

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/BARGHEST-ngo/MESH/provisioning/internal/state"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

// Todo - pin version?
const frps_image = "snowdreamtech/frps:latest"

func Start(r *state.Registry, slug string) error {
	d, ok := r.Get(slug)
	if !ok {
		return fmt.Errorf("deployment not defined")
	}

	configPath, err := writeConfig(d)
	if err != nil {
		return fmt.Errorf("failed to write frps config file: %w", err)
	}

	client, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return fmt.Errorf("failed to create docker client: %w", err)
	}
	defer client.Close()

	ctx := context.Background()
	output, err := client.ImagePull(ctx, frps_image, image.PullOptions{})
	if err != nil {
		return fmt.Errorf("failed to pull frps image: %w", err)
	}
	io.Copy(io.Discard, output)
	output.Close()
	resp, err := client.ContainerCreate(ctx,
		&container.Config{
			Image: frps_image,
		},
		&container.HostConfig{
			RestartPolicy: container.RestartPolicy{Name: container.RestartPolicyUnlessStopped},
			Binds:         []string{fmt.Sprintf("%s:/etc/frp/frps.toml:ro", configPath)},
			PortBindings:  nat.PortMap{nat.Port("7000/tcp"): []nat.PortBinding{{HostPort: fmt.Sprintf("%d", d.FrpsPort)}}},
		},
		&network.NetworkingConfig{},
		nil, fmt.Sprintf("frps-%s", slug))
	if err != nil {
		return fmt.Errorf("failed to create container: %w", err)
	}

	fmt.Printf("container id: %s", resp.ID)

	return nil
}

func writeConfig(d state.Deployment) (string, error) {
	dir := filepath.Join(os.TempDir(), d.Slug)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return "", err
	}

	path := filepath.Join(dir, "frps.toml")
	content := fmt.Sprintf("bindPort = 7000\nauth.token = %q\n", d.Token)
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		return "", err
	}
	return path, nil
}
