package container

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"

	"github.com/synalice/gobox/docker/controller"
)

// Start starts a container
func Start(controller *controller.Controller, containerID string) error {
	err := controller.Cli.ContainerStart(context.Background(), containerID, types.ContainerStartOptions{})
	if err != nil {
		return fmt.Errorf("error starting container: %w", err)
	}

	return nil
}
