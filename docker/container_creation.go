// TODO: Add time limit setting for each new container
// TODO: Add graceful shutdown

package docker

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/google/uuid"

	v1 "github.com/opencontainers/image-spec/specs-go/v1"
)

// NewController returns controller that wll be used for running methods off of it
func NewController() (c *Controller, err error) {
	c = new(Controller)

	c.cli, err = client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, fmt.Errorf("error while executing client.NewClientWithOpts(): %w", err)
	}

	return c, nil
}

// EnsureImage pulls images from docker hub to make sure in exists
func (c *Controller) EnsureImage(config *ContainerConfig) error {
	if config.LocalImage {
		return nil
	}

	reader, err := c.cli.ImagePull(context.Background(), config.Image, types.ImagePullOptions{})
	if err != nil {
		return fmt.Errorf("error while executing c.cli.ImagePull(): %w", err)
	}
	defer reader.Close()

	_, err = io.Copy(os.Stdout, reader)
	if err != nil {
		return fmt.Errorf("error coping reader to stdout: %w", err)
	}

	return nil
}

// ContainerCreate creates a container
func (c *Controller) ContainerCreate(config *ContainerConfig, volumes []VolumeMount) (id string, err error) {
	hostConfig := container.HostConfig{}
	networkingConfig := network.NetworkingConfig{}
	platform := v1.Platform{}

	var mounts []mount.Mount

	for _, volume := range volumes {
		m := mount.Mount{
			Type:   mount.TypeVolume,
			Source: volume.Volume.Name,
			Target: volume.HostPath,
		}
		mounts = append(mounts, m)
	}

	hostConfig.Mounts = mounts

	resp, err := c.cli.ContainerCreate(
		context.Background(),
		&container.Config{
			Tty:   true,
			Image: config.Image,
			Cmd:   config.Cmd,
		},
		&hostConfig,
		&networkingConfig,
		&platform,
		generateContainerName(config),
	)
	if err != nil {
		return "", fmt.Errorf("error while executing c.cli.ContainerCreate(): %w", err)
	}

	return resp.ID, nil
}

// ContainerStart starts a container
func (c *Controller) ContainerStart(containerID string) error {
	err := c.cli.ContainerStart(context.Background(), containerID, types.ContainerStartOptions{})
	if err != nil {
		return fmt.Errorf("error while executing c.cli.ContainerStart(): %w", err)
	}

	return nil
}

// ContainerWait waits until the container is stopped
func (c *Controller) ContainerWait(id string) (state int64, err error) {
	resultC, errC := c.cli.ContainerWait(context.Background(), id, "not-running")
	select {
	case err := <-errC:
		return 0, err
	case result := <-resultC:
		return result.StatusCode, nil
	}
}

// ContainerLog returns logs of a specific container
func (c *Controller) ContainerLog(id string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	reader, err := c.cli.ContainerLogs(ctx, id, types.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
	})
	if err != nil {
		return "", fmt.Errorf("error while executing c.cli.ContainerLogs(): %w", err)
	}

	buffer, err := io.ReadAll(reader)
	if err != nil && err != io.EOF {
		return "", fmt.Errorf("error while executing io.ReadAll(): %w", err)
	}

	return string(buffer), nil
}

// ContainerRemove removes a container
func (c *Controller) ContainerRemove(id string) error {
	err := c.cli.ContainerRemove(context.Background(), id, types.ContainerRemoveOptions{
		RemoveVolumes: true,
		Force:         true,
	})
	if err != nil {
		return fmt.Errorf("error executing c.cli.ContainerRemove() %w", err)
	}

	return nil
}

// Run creates, runs and removes container defined by the ContainerConfig
func (c *Controller) Run(config *ContainerConfig, volumes []VolumeMount) (statusCode int64, body string, err error) {
	// TODO: Create custom errors for this

	// Pulls image if needed
	err = c.EnsureImage(config)
	if err != nil {
		return statusCode, body, err
	}

	// Create the container
	id, err := c.ContainerCreate(config, volumes)
	if err != nil {
		return statusCode, body, err
	}

	// Run the container
	err = c.ContainerStart(id)
	if err != nil {
		return statusCode, body, err
	}

	// Wait for it to finish
	statusCode, err = c.ContainerWait(id)
	if err != nil {
		return statusCode, body, err
	}

	// Get the log
	body, err = c.ContainerLog(id)
	if err != nil {
		return statusCode, body, err
	}

	// Remove the container
	err = c.cli.ContainerRemove(context.Background(), id, types.ContainerRemoveOptions{
		RemoveVolumes: true,
		Force:         true,
	})
	if err != nil {
		return statusCode, body, fmt.Errorf("unable to remove container %q: %w", id, err)
	}

	return statusCode, body, err
}

// generateContainerName generates unique name for each new container
func generateContainerName(config *ContainerConfig) (containerName string) {
	return "gobox" + "-" + config.Image + "-" + uuid.NewString()
}
