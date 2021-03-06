package volume

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	volumetypes "github.com/docker/docker/api/types/volume"
	"github.com/google/uuid"

	"github.com/synalice/gobox/docker/controller"
)

// Find finds a specified volume
func Find(controller *controller.Controller, volumeName string) (volume *types.Volume, err error) {
	volumes, err := controller.Cli.VolumeList(context.Background(), filters.NewArgs())
	if err != nil {
		return nil, err
	}

	for _, v := range volumes.Volumes {
		if v.Name == volumeName {
			return v, nil
		}
	}

	return nil, nil
}

// Ensure makes sure specified volume exists and creates it if it is not
// present. Use empty string to generate name randomly.
func Ensure(controller *controller.Controller, volumeName string) (volume *types.Volume, err error) {
	if volumeName == "" {
		vol, err := controller.Cli.VolumeCreate(context.Background(), volumetypes.VolumeCreateBody{
			Driver: "local",
			Name:   generateVolumeName(),
		})
		return &vol, err
	}

	volume, err = Find(controller, volumeName)
	if err != nil {
		return nil, err
	}
	if volume != nil {
		return volume, nil
	}

	vol, err := controller.Cli.VolumeCreate(context.Background(), volumetypes.VolumeCreateBody{
		Driver: "local",
		Name:   volumeName,
	})

	return &vol, err
}

// Remove removes specified volume
func Remove(controller *controller.Controller, volumeName string) error {
	vol, err := Find(controller, volumeName)
	if err != nil {
		return fmt.Errorf("couldn't find volume: %w", err)
	}
	if vol == nil {
		return nil
	}

	err = controller.Cli.VolumeRemove(context.Background(), volumeName, true)
	if err != nil {
		return fmt.Errorf("couldn't remove volume: %w", err)
	}

	return nil
}

// generateVolumeName generates unique UUID name for each new volume
func generateVolumeName() string {
	return "gobox" + "-" + "volume" + "-" + uuid.NewString()
}
