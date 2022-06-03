package controller

import (
	"fmt"

	"github.com/docker/docker/client"
)

// Controller is an object required for all operations on docker
type Controller struct {
	Cli *client.Client
}

// NewController returns new Controller that wll be used for running methods
// off of it
func NewController() (c *Controller, err error) {
	c = &Controller{
		Cli: nil,
	}

	c.Cli, err = client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, fmt.Errorf("error creating new Controller: %w", err)
	}

	return c, nil
}
