package container

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/synalice/gobox/docker/config"
	"github.com/synalice/gobox/docker/controller"
)

type Builder struct {
	controller *controller.Controller
	config     *config.Config
	container  *Container
}

func NewContainerBuilder(controller *controller.Controller) *Builder {
	return &Builder{
		controller: controller,
		config:     nil,
		container:  &Container{},
	}
}

func (b *Builder) SetConfig(config *config.Config) *Builder {
	b.config = config
	return b
}

func (b *Builder) Build() (*Container, error) {
	err := b.createContainer()
	if err != nil {
		return nil, fmt.Errorf("couldn't create container: %w", err)
	}

	b.setTimeLimit(b.config.TimeLimit)

	return b.container, nil
}

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

func (b *Builder) setTimeLimit(timeLimit time.Duration) *Builder {
	b.container.TimeLimit = timeLimit
	return b
}

func (b *Builder) generateUUIDName() (containerName string) {
	return "gobox" + "-" + b.config.ContainerConfig.Image + "-" + uuid.NewString()
}
