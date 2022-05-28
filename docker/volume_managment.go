package docker

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"

	volumetypes "github.com/docker/docker/api/types/volume"
)

// FindVolume finds a specified volume
func (c *Controller) FindVolume(name string) (volume *types.Volume, err error) {
	volumes, err := c.Cli.VolumeList(context.Background(), filters.NewArgs())
	if err != nil {
		return nil, err
	}

	for _, v := range volumes.Volumes {
		if v.Name == name {
			return v, nil
		}
	}

	return nil, nil
}

// EnsureVolume makes sure specified volume exists and creates it if it doesn't
func (c *Controller) EnsureVolume(name string) (new bool, volume *types.Volume, err error) {
	volume, err = c.FindVolume(name)
	if err != nil {
		return false, nil, err
	}
	if volume != nil {
		return false, volume, nil
	}

	vol, err := c.Cli.VolumeCreate(context.Background(), volumetypes.VolumeCreateBody{
		Driver: "local",
		Name:   name,
	})

	return true, &vol, err
}

// RemoveVolume removes specified volume
func (c *Controller) RemoveVolume(name string) error {
	vol, err := c.FindVolume(name)
	if err != nil {
		return fmt.Errorf("couldn't find volume: %w", err)
	}
	if vol == nil {
		return nil
	}

	err = c.Cli.VolumeRemove(context.Background(), name, true)
	if err != nil {
		return fmt.Errorf("couldn't remove volume: %w", err)
	}

	return nil
}
