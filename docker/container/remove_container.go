package container

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"

	"github.com/synalice/gobox/docker/controller"
)

// Remove deletes a specific container from host
func Remove(controller *controller.Controller, containerIDs ...string) error {
	for _, containerID := range containerIDs {
		err := controller.Cli.ContainerRemove(context.Background(), containerID, types.ContainerRemoveOptions{
			RemoveVolumes: true,
			Force:         true,
		})
		if err != nil {
			return fmt.Errorf("error removing container: %w", err)
		}
	}

	return nil
}
