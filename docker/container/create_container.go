package container

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/synalice/gobox/docker/config"
	"github.com/synalice/gobox/docker/controller"
)

// Builder is a structure used for building new docker container
type Builder struct {
	controller *controller.Controller
	config     *config.Config
	container  *Container
}

// NewContainerBuilder returns a Builder object that will be used for building
// a docker container
func NewContainerBuilder(controller *controller.Controller) *Builder {
	return &Builder{
		controller: controller,
		config:     nil,
		container:  &Container{},
	}
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

	b.setTimeLimit(b.config.TimeLimit)

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
func (b *Builder) setTimeLimit(timeLimit time.Duration) *Builder {
	b.container.TimeLimit = timeLimit
	return b
}

// generateUUIDName generates unique UUID name for each new container
func (b *Builder) generateUUIDName() (containerName string) {
	return "gobox" + "-" + b.config.ContainerConfig.Image + "-" + uuid.NewString()
}
