package container

import (
	"bytes"
	"fmt"
	"io"

	"github.com/synalice/gobox/docker/controller"
)

// GetLogs returns logs of a specific container
func GetLogs(controller *controller.Controller, container *Container) (string, error) {
	// GetLogs doesn't really need a controller, but I think it's better to
	// still require one to make the user API more unified.
	//
	// This line is given here only so that the compiler doesn't yell at
	// you because of an unused parameter.
	_ = controller

	r := container.Connection.Reader
	var buf bytes.Buffer

	_, err := io.Copy(&buf, r)
	if err != nil {
		return "", fmt.Errorf("error getting contaner's logs: %w", err)
	}

	return buf.String(), nil
}
