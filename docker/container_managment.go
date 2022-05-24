// TODO: Add time limit setting for each new container
// TODO: Add graceful shutdown

package docker

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/google/uuid"

	v1 "github.com/opencontainers/image-spec/specs-go/v1"
)

// NewController returns Controller that wll be used for running methods off of it
func NewController(config *ContainerConfig) (c *Controller, err error) {
	c = &Controller{
		cli:    nil,
		config: config,
	}

	c.cli, err = client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, fmt.Errorf("error while executing client.NewClientWithOpts(): %w", err)
	}

	return c, nil
}

// EnsureImage pulls images from docker hub to make sure in exists
func (c *Controller) EnsureImage(imageName string) error {
	if c.config.LocalImage {
		return nil
	}

	reader, err := c.cli.ImagePull(context.Background(), imageName, types.ImagePullOptions{})
	if err != nil {
		return fmt.Errorf("error while executing c.cli.ImagePull(): %w", err)
	}
	defer reader.Close()

	// FIXME: These 4 lines are unused
	//_, err = io.Copy(os.Stdout, reader)
	//if err != nil {
	//	return fmt.Errorf("error coping reader to stdout: %w", err)
	//}

	return nil
}

// ContainerCreate creates a container
func (c *Controller) ContainerCreate(config *ContainerConfig, volumes []VolumeMount) (containerID string, err error) {
	hostConfig := container.HostConfig{
		Resources: container.Resources{
			Memory: config.MemoryMB * 1024 * 1024,
		},
	}
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
func (c *controller) ContainerWait(containerID string) (state int64, err error) {
	resultC, errC := c.cli.ContainerWait(context.Background(), containerID, "not-running")
	select {
	case err := <-errC:
		return 0, err
	case result := <-resultC:
		return result.StatusCode, nil
	}
}

// ContainerLog returns logs of a specific container
func (c *Controller) ContainerLog(containerID string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	reader, err := c.cli.ContainerLogs(ctx, containerID, types.ContainerLogsOptions{
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
func (c *Controller) ContainerRemove(containerID string) error {
	err := c.cli.ContainerRemove(context.Background(), containerID, types.ContainerRemoveOptions{
		RemoveVolumes: true,
		Force:         true,
	})
	if err != nil {
		return fmt.Errorf("error executing c.cli.ContainerRemove() %w", err)
	}

	return nil
}

// Run creates, runs and removes container defined by the ContainerConfig
func (c *Controller) Run(volumes []VolumeMount) (statusCode int64, logs string, err error) {
	// TODO: Create custom errors for this

	// Pulls image if needed
	err = c.EnsureImage(c.config.Image)
	if err != nil {
		return statusCode, logs, err
	}

	// Create the container
	id, err := c.ContainerCreate(c.config, volumes)
	if err != nil {
		return statusCode, logs, err
	}

	// Start the container
	err = c.ContainerStart(id)
	if err != nil {
		return statusCode, logs, err
	}

	// Wait for it to finish
	statusCode, err = c.ContainerWait(id, c.config.TimeLimit)
	if err != nil {
		return statusCode, logs, err
	}

	// Get the log
	logs, err = c.ContainerLog(id)
	if err != nil {
		return statusCode, logs, err
	}

	// Remove the container
	err = c.cli.ContainerRemove(context.Background(), id, types.ContainerRemoveOptions{
		RemoveVolumes: true,
		Force:         true,
	})
	if err != nil {
		return statusCode, logs, fmt.Errorf("unable to remove container %q: %w", id, err)
	}

	return statusCode, logs, err
}

// generateContainerName generates unique name for each new container
func generateContainerName(config *ContainerConfig) (containerName string) {
	return "gobox" + "-" + config.Image + "-" + uuid.NewString()
}
