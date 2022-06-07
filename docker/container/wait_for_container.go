package container

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/synalice/gobox/docker/controller"
)

var (
	ErrorTimeout = errors.New("container killed due to timeout")
)

// Wait waits until the container has either finished executing or was killed
// for exceeding time limit
func Wait(controller *controller.Controller, containerID string, timeLimit time.Duration) (state int64, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeLimit)
	defer cancel()

	resultC, errC := controller.Cli.ContainerWait(ctx, containerID, "")
	select {
	case err := <-errC:
		if err.Error() == "context deadline exceeded" {
			return 0, ErrorTimeout
		}
		return 0, fmt.Errorf("error while waiting for container to finish: %w", err)
	case result := <-resultC:
		return result.StatusCode, nil
	}
}
