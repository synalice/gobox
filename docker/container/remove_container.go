package container

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"

	"github.com/synalice/gobox/docker/controller"
)

// Remove deletes a container
func Remove(controller *controller.Controller, containerID string) error {
	err := controller.Cli.ContainerRemove(context.Background(), containerID, types.ContainerRemoveOptions{
		RemoveVolumes: true,
		Force:         true,
	})
	if err != nil {
		return fmt.Errorf("error removing container: %w", err)
	}

	return nil
}
