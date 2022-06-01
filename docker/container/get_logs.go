package container

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/docker/docker/api/types"

	"github.com/synalice/gobox/docker/controller"
)

// GetLogs returns logs of a specific container
func GetLogs(controller *controller.Controller, containerID string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	reader, err := controller.Cli.ContainerLogs(ctx, containerID, types.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
	})
	if err != nil {
		return "", fmt.Errorf("error getting contaner's logs: %w", err)
	}

	buffer, err := io.ReadAll(reader)
	if err != nil && err != io.EOF {
		return "", fmt.Errorf("error while executing io.ReadAll(): %w", err)
	}

	return string(buffer), nil
}
