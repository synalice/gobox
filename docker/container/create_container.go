package container

import (
	"context"
	"fmt"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/google/uuid"

	"github.com/synalice/gobox/docker/config"
	"github.com/synalice/gobox/docker/controller"
	"github.com/synalice/gobox/docker/file"
	"github.com/synalice/gobox/docker/mount"
)

type fileInContainer struct {
	file         file.File
	fileLocation string
}

// Builder is a structure used for building new docker container
type Builder struct {
	controller *controller.Controller
	config     *config.Config
	files      []fileInContainer
	container  *Container
}

// NewContainerBuilder returns a Builder object that will be used for building
// a docker container
func NewContainerBuilder(controller *controller.Controller) *Builder {
	return &Builder{
		controller: controller,
		config:     nil,
		container:  &Container{},
		files:      nil,
	}
}

// SetFile creates files that will be executed inside the container. Can be
// called multiple times for multiple mounts.
func (b *Builder) SetFile(file file.File, location mount.Mount) *Builder {
	fileInContainer := fileInContainer{
		file:         file,
		fileLocation: location.Folder,
	}

	b.files = append(b.files, fileInContainer)

	return b
}

// SetConfig sets config that will be later used for building a docker container
func (b *Builder) SetConfig(config *config.Config) *Builder {
	b.config = config
	return b
}

// Build builds a new docker container
func (b *Builder) Build() (*Container, error) {
	err := b.createContainer()
	if err != nil {
		return nil, fmt.Errorf("couldn't create container: %w", err)
	}

	err = b.copyFilesToContainer()
	if err != nil {
		return nil, fmt.Errorf("couldn't create container: %w", err)
	}

	b.setTimeLimit()

	return b.container, nil
}

// createContainer works directly with docker API and creates new container
func (b *Builder) createContainer() error {
	resp, err := b.controller.Cli.ContainerCreate(
		context.Background(),
		&b.config.ContainerConfig,
		&b.config.HostConfig,
		&b.config.NetworkingConfig,
		&b.config.Platform,
		b.generateUUIDName(),
	)
	if err != nil {
		return fmt.Errorf("error building container: %w", err)
	}

	b.container.ID = resp.ID

	return nil
}

// setTimeLimit sets maximum allowed time for the container to run
func (b *Builder) setTimeLimit() *Builder {
	b.container.TimeLimit = b.config.TimeLimit
	return b
}

// generateUUIDName generates unique UUID name for each new container
func (b *Builder) generateUUIDName() (containerName string) {
	if strings.Contains(b.config.ContainerConfig.Image, ":") {
		newImgName := strings.Replace(b.config.ContainerConfig.Image, ":", "-", -1)
		return "gobox" + "-" + newImgName + "-" + uuid.NewString()
	}
	return "gobox" + "-" + b.config.ContainerConfig.Image + "-" + uuid.NewString()
}

// copyFilesToContainer iterates over Builder.files and copies each file into
// the container
func (b *Builder) copyFilesToContainer() error {
	for _, f := range b.files {
		archive, err := file.CreateTarball([]file.File{f.file})
		if err != nil {
			return fmt.Errorf("error copying files into the container: %w", err)
		}

		err = b.controller.Cli.CopyToContainer(
			context.Background(),
			b.container.ID,
			f.fileLocation,
			&archive,
			types.CopyToContainerOptions{},
		)
		if err != nil {
			return fmt.Errorf("error copying files into the container: %w", err)
		}
	}

	return nil
}
