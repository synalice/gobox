package config

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"

	"github.com/docker/docker/api/types"

	"github.com/synalice/gobox/docker/controller"
)

// EnsureImage pulls an image from the docker hub to make sure in exists
func EnsureImage(controller *controller.Controller, imageName string) error {
	reader, err := controller.Cli.ImagePull(context.Background(), imageName, types.ImagePullOptions{})
	if err != nil {
		return fmt.Errorf("error while pulling an image from a remote registry: %w", err)
	}
	defer func(reader io.ReadCloser) {
		err := reader.Close()
		if err != nil {
			log.Println("error while pulling an image from a remote registry:", err)
		}
	}(reader)

	_, err = io.Copy(ioutil.Discard, reader)
	if err != nil {
		return fmt.Errorf("error while pulling an image from a remote registry: %w", err)
	}

	return nil
}
