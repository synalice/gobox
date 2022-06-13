package container

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"

	"github.com/synalice/gobox/docker/controller"
)

// Start starts a container
func Start(controller *controller.Controller, containerID string, stdin string) error {
	if stdin != "" {
		err := startWithStdin(controller, containerID, stdin)
		if err != nil {
			return fmt.Errorf("error starting container: %w", err)
		}

		return nil
	}

	err := startWithoutStdin(controller, containerID)
	if err != nil {
		return fmt.Errorf("error starting container: %w", err)
	}

	return nil
}

func startWithStdin(controller *controller.Controller, containerID string, stdin string) error {
	hijackedResponse, err := attachToContainer(controller, containerID)
	if err != nil {
		return err
	}

	err = controller.Cli.ContainerStart(context.Background(), containerID, types.ContainerStartOptions{})
	if err != nil {
		return err
	}

	err = writeToStdin(hijackedResponse, stdin)
	if err != nil {
		return err
	}

	return nil
}

func startWithoutStdin(controller *controller.Controller, containerID string) error {
	err := controller.Cli.ContainerStart(context.Background(), containerID, types.ContainerStartOptions{})
	if err != nil {
		return err
	}

	return nil
}

func attachToContainer(controller *controller.Controller, containerID string) (types.HijackedResponse, error) {
	hijackedResponse, err := controller.Cli.ContainerAttach(
		context.Background(),
		containerID,
		types.ContainerAttachOptions{
			Stream: true,
			Stdin:  true,
			Stdout: true,
			Stderr: true,
		},
	)
	if err != nil {
		return hijackedResponse, fmt.Errorf("error attaching to a container: %w", err)
	}

	return hijackedResponse, nil
}

func writeToStdin(connection types.HijackedResponse, input string) error {
	_, err := connection.Conn.Write([]byte(input + "\n"))
	if err != nil {
		return fmt.Errorf("error writing into stdin: %w", err)
	}

	return nil
}
