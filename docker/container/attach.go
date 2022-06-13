package container

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"

	"github.com/synalice/gobox/docker/controller"
)

// attach is basically an equivalent of `$ docker attach <container>`
// It returns a connection to a container, allowing you to read from stdout and
// write into stdin.
func attach(controller *controller.Controller, containerID string) (types.HijackedResponse, error) {
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
