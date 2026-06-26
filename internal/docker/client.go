package docker

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	dockerclient "github.com/docker/docker/client"
)

// Client wraps the official Docker SDK client.
type Client struct {
	dc *dockerclient.Client
}

// NewClient creates a Client connected to the local Docker daemon via DOCKER_HOST / Unix socket.
func NewClient() (*Client, error) {
	dc, err := dockerclient.NewClientWithOpts(
		dockerclient.FromEnv,
		dockerclient.WithAPIVersionNegotiation(),
	)
	if err != nil {
		return nil, err
	}
	return &Client{dc: dc}, nil
}

// Close releases the underlying Docker client.
func (c *Client) Close() error { return c.dc.Close() }

// FetchStates returns a ContainerState for each service name.
// Services not found in the running container list are reported as StatePending.
func (c *Client) FetchStates(ctx context.Context, services []string) (map[string]ContainerState, error) {
	f := filters.NewArgs()
	f.Add("label", "com.docker.compose.service")

	containers, err := c.dc.ContainerList(ctx, container.ListOptions{
		All:     true,
		Filters: f,
	})
	if err != nil {
		return nil, err
	}

	states := make(map[string]ContainerState, len(services))
	for _, svc := range services {
		states[svc] = StatePending
	}

	for _, ctr := range containers {
		svcName := ctr.Labels["com.docker.compose.service"]
		if _, ok := states[svcName]; !ok {
			continue
		}
		states[svcName] = mapContainerState(ctr)
	}
	return states, nil
}

// mapContainerState converts a Docker container summary to a ContainerState.
// ctr.Status is a human-readable string like "Up 2 minutes (healthy)" — we
// parse it with strings.Contains rather than an exact switch.
func mapContainerState(ctr container.Summary) ContainerState {
	switch ctr.State {
	case "running":
		st := ctr.Status
		switch {
		case strings.Contains(st, "(healthy)"):
			return StateHealthy
		case strings.Contains(st, "(health: starting)"):
			return StateStarting
		case strings.Contains(st, "(unhealthy)"):
			return StateUnhealthy
		default:
			// Running with no healthcheck — treat as healthy.
			return StateHealthy
		}
	case "exited":
		return StateExited
	default:
		return StateStarting
	}
}

// Logs returns the last n lines of stdout+stderr for containerID.
func (c *Client) Logs(ctx context.Context, containerID string, lines int) (string, error) {
	tail := "200"
	if lines > 0 {
		tail = fmt.Sprintf("%d", lines)
	}
	opts := container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Tail:       tail,
	}
	reader, err := c.dc.ContainerLogs(ctx, containerID, opts)
	if err != nil {
		return "", err
	}
	defer func() { _ = reader.Close() }()

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, reader); err != nil {
		return "", err
	}
	return buf.String(), nil
}
