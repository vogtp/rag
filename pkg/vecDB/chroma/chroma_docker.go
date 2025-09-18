package chroma

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/spf13/viper"
	"github.com/vogtp/rag/pkg/cfg"
)

type Container struct {
	slog *slog.Logger

	containerName string
	imageName     string

	cli *client.Client
}

func NewContainer(slog *slog.Logger) (*Container, error) {
	containerName := viper.GetString(cfg.ChromaContainerName)
	imageName := viper.GetString(cfg.ChromaContainerImage)
	c := &Container{
		containerName: containerName,
		imageName:     imageName,
		slog:          slog.With("containerName", containerName, "imageName", imageName),
	}
	cli, err := client.NewClientWithOpts(
		client.FromEnv,
		client.WithAPIVersionNegotiation(),
	)
	if err != nil {
		return nil, fmt.Errorf("cannot create docker client: %w", err)
	}
	c.cli = cli
	return c, nil
}

func EnsureStarted(slog *slog.Logger, ctx context.Context) error {
	c, err := NewContainer(slog)
	if err != nil {
		return err
	}
	_, err = c.EnsureStarted(ctx)
	return err
}

func (c *Container) EnsureStarted(ctx context.Context) (func(ctx context.Context) error, error) {

	id, running, err := c.getOrCreateContainer(ctx)
	if err != nil {
		return nil, fmt.Errorf("cannot find docker image: %w", err)
	}
	stopContainer := func(ctx context.Context) error {
		return c.stopContainer(ctx, id)
	}
	if running {
		c.slog.Info("Container already running")
		return stopContainer, nil
	}
	c.slog.Info("Starting container")
	if err := c.cli.ContainerStart(ctx, id, container.StartOptions{}); err != nil {
		return nil, fmt.Errorf("cannot start docker container %s: %w", c.containerName, err)
	}
	return stopContainer, nil
}

func (c *Container) stopContainer(ctx context.Context, id string) error {
	c.slog.Info("Stopping container")
	if err := c.cli.ContainerStop(ctx, id, container.StopOptions{}); err != nil {
		return fmt.Errorf("cannot stop docker container %s: %w", c.containerName, err)
	}
	return nil
}

func (c *Container) EnsureStopped(ctx context.Context) error {
	id, running, err := c.getOrCreateContainer(ctx)
	if err != nil {
		return fmt.Errorf("cannot find docker image: %w", err)
	}
	if !running {
		return nil
	}
	return c.stopContainer(ctx, id)
}

func (c *Container) getOrCreateContainer(ctx context.Context) (string, bool, error) {
	containers, err := c.cli.ContainerList(ctx, container.ListOptions{All: true})
	if err != nil {
		return "", false, err
	}
	for _, container := range containers {
		c.slog.Debug("searching container", "foundImage", container.Image, "state", container.State)
		if container.Image == c.imageName {
			c.slog.Debug("Found existing container", "state", container.State)
			state := container.State == "running"
			return container.ID, state, nil
		}
	}
	id, err := c.create(ctx)
	return id, false, err
}

func (c *Container) create(ctx context.Context) (string, error) {
	contPort := "8000/tcp"
	port := viper.GetString(cfg.ChromaContainerPort)
	containerCfg := &container.Config{
		Image:        c.imageName,
		AttachStdout: true,
		AttachStderr: true,
		ExposedPorts: map[nat.Port]struct{}{
			nat.Port(contPort): {},
		},
		Env: []string{"IS_PERSISTENT=TRUE", "ANONYMIZED_TELEMETRY=FALSE", fmt.Sprintf("CHROMA_HOST_PORT=%s", "8000")},
	}

	// Set port binding
	hostBinding := nat.PortBinding{
		HostIP:   "0.0.0.0",
		HostPort: port,
	}

	containerPortBinding := nat.PortMap{
		nat.Port(contPort): []nat.PortBinding{hostBinding},
	}

	// Create host configuration
	hostCfg := &container.HostConfig{
		PortBindings: containerPortBinding,
	}

	c.slog.Info("Pulling image")
	out, err := c.cli.ImagePull(ctx, c.imageName, image.PullOptions{})
	if err != nil {
		c.slog.Warn("docker pull", "err", err)
		return "", fmt.Errorf("cannot pull image %s: %w", c.imageName, err)
	}
	defer out.Close()
	s := bufio.NewScanner(out)
	for s.Scan() {
		var resp map[string]any
		if err := json.Unmarshal(s.Bytes(), &resp); err != nil {
			slog.Warn("Cannot unmarshal docker pull response", "err", err, "response", s.Text())
		}
		sl := c.slog

		for k, v := range resp {
			if k == "status" {
				continue
			}
			sl = sl.With(slog.Any(k, v))
		}
		sl.Info(fmt.Sprintf("%v", resp["status"]))
	}
	// hostCfg := &container.HostConfig{
	// 	PortBindings: map[nat.Port][]nat.PortBinding{nat.Port(contPort): {{HostPort: port}}},

	// 	// Mounts: []mount.Mount{
	// 	// 	{
	// 	// 		Type:   mount.TypeVolume,
	// 	// 		Target: "./chroma",
	// 	// 		Source: "/chroma/chroma",
	// 	// 	},
	// 	// },
	// 	// Binds: []string{
	// 	// 	"./chroma:/srv/chroma",
	// 	// },
	// }
	resp, err := c.cli.ContainerCreate(ctx, containerCfg, hostCfg, nil, nil, c.containerName)
	if err != nil {
		c.slog.Warn("docker create", "err", err)
		return "", fmt.Errorf("cannot prepaire docker container %s: %w", c.containerName, err)
	}
	return resp.ID, nil
}
