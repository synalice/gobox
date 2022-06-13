package container

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"

	"github.com/synalice/gobox/docker/controller"
)

// Start starts a container
func Start(controller *controller.Controller, container *Container, stdin string) error {
	if stdin != "" {
		err := startWithStdin(controller, container, stdin)
		if err != nil {
			return fmt.Errorf("error starting container: %w", err)
		}

		return nil
	}

	err := startWithoutStdin(controller, container)
	if err != nil {
		return fmt.Errorf("error starting container: %w", err)
	}

	return nil
}

func startWithStdin(controller *controller.Controller, container *Container, stdin string) error {
	err := controller.Cli.ContainerStart(context.Background(), container.ID, types.ContainerStartOptions{})
	if err != nil {
		return err
	}

	err = writeToStdin(container.Connection, stdin)
	if err != nil {
		return err
	}

	return nil
}

func startWithoutStdin(controller *controller.Controller, container *Container) error {
	err := controller.Cli.ContainerStart(context.Background(), container.ID, types.ContainerStartOptions{})
	if err != nil {
		return err
	}

	return nil
}

func writeToStdin(connection types.HijackedResponse, input string) error {
	_, err := connection.Conn.Write([]byte(input + "\n"))
	if err != nil {
		return fmt.Errorf("error writing into stdin: %w", err)
	}

	return nil
}
