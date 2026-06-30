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

type Manager struct {
	FrpsImage string
}

func (m Manager) Start(d state.Deployment) error {
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
	output, err := client.ImagePull(ctx, m.FrpsImage, image.PullOptions{})
	if err != nil {
		return fmt.Errorf("failed to pull frps image: %w", err)
	}
	io.Copy(io.Discard, output)
	output.Close()
	resp, err := client.ContainerCreate(ctx,
		&container.Config{
			Image: m.FrpsImage,
			Labels: map[string]string{
				"traefik.enable": "true",
				fmt.Sprintf("traefik.http.routers.%s.rule", d.Slug):                      fmt.Sprintf("Host(`%s.tunnel.meshforensics.app`)", d.Slug),
				fmt.Sprintf("traefik.http.routers.%s.tls", d.Slug):                       "true",
				fmt.Sprintf("traefik.http.services.%s.loadbalancer.server.port", d.Slug): "8080",
				"traefik.docker.network":                                                 "mesh-proxy",
			},
		},
		&container.HostConfig{
			RestartPolicy: container.RestartPolicy{Name: container.RestartPolicyUnlessStopped},
			Binds:         []string{fmt.Sprintf("%s:/etc/frp/frps.toml:ro", configPath)},
			PortBindings: nat.PortMap{
				nat.Port("7000/tcp"): []nat.PortBinding{{HostPort: fmt.Sprintf("%d", d.FrpsPort)}},
			},
		},
		&network.NetworkingConfig{
			EndpointsConfig: map[string]*network.EndpointSettings{
				"mesh-proxy": {},
			},
		},
		nil, fmt.Sprintf("frps-%s", d.Slug))
	if err != nil {
		return fmt.Errorf("failed to create container: %w", err)
	}

	return client.ContainerStart(ctx, resp.ID, container.StartOptions{})
}

func writeConfig(d state.Deployment) (string, error) {
	hostDataPath := os.Getenv("HOST_DATA_PATH")
	if hostDataPath == "" {
		return "", fmt.Errorf("HOST_DATA_PATH not set")
	}

	dir := filepath.Join(hostDataPath, d.Slug)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return "", err
	}

	path := filepath.Join(dir, "frps.toml")
	content := fmt.Sprintf("bindPort = 7000\nauth.token = %q\nvhostHTTPPort = 8080\n", d.Token)
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		return "", err
	}
	return path, nil
}

func (Manager) Stop(slug string) error {
	c, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}
	defer c.Close()

	ctx := context.Background()
	name := fmt.Sprintf("frps-%s", slug)
	if err := c.ContainerStop(ctx, name, container.StopOptions{}); err != nil {
		return err
	}

	return c.ContainerRemove(ctx, name, container.RemoveOptions{})
}
