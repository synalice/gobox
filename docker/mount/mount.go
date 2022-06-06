package mount

import (
	"fmt"

	"github.com/docker/docker/api/types"

	"github.com/synalice/gobox/docker/controller"
	"github.com/synalice/gobox/docker/volume"
)

type Mount struct {
	Volume *types.Volume
	Folder string
}

// NewMount returns new Mount object that can be later used to specify a mount
// point when building a config.
func NewMount(controller *controller.Controller, volumeName string, folderInContainer string) (Mount, error) {
	newVolume, err := volume.Ensure(controller, volumeName)
	if err != nil {
		return Mount{}, fmt.Errorf("error creating new mount: %w", err)
	}

	err = validatePathToFolder(folderInContainer)
	if err != nil {
		return Mount{}, err
	}

	mount := Mount{
		Volume: newVolume,
		Folder: folderInContainer,
	}

	return mount, nil
}

// validatePathToFolder makes sure that path to the folder is absolute
func validatePathToFolder(path string) error {
	if path[0:1] != "/" {
		return fmt.Errorf("invalid mount path: '%v' mount path must be absolute", path)
	}

	return nil
}

// CleanUp takes in an array of mounts and removes all volumes in them.
// Use this function after you've finished all the work and want to remove all
// leftover volumes.
func CleanUp(controller *controller.Controller, mounts ...Mount) error {
	for _, mount := range mounts {
		err := volume.Remove(controller, mount.Volume.Name)
		if err != nil {
			return fmt.Errorf("error removing used mounts: %w", err)
		}
	}

	return nil
}
