package container

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/docker/docker/api/types"

	"github.com/synalice/gobox/docker/controller"
)

// EnsureImage pulls images from docker hub to make sure in exists
func EnsureImage(controller *controller.Controller, imageName string) error {
	reader, err := controller.Cli.ImagePull(context.Background(), imageName, types.ImagePullOptions{})
	if err != nil {
		return fmt.Errorf("error while executing c.cli.ImagePull(): %w", err)
	}
	defer func(reader io.ReadCloser) {
		err := reader.Close()
		if err != nil {
			log.Fatalf("couldn't close the reader: %v", err)
		}
	}(reader)

	return nil
}
