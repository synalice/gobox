package container

import (
	"time"

	"github.com/docker/docker/api/types"
)

// Container struct holds information about newly created container
type Container struct {
	// ID of the container
	ID string

	// Maximum allowed time for the container to run
	TimeLimit time.Duration

	// A connection to a container that allows you to read from stdout and
	// write into stdin.
	Connection types.HijackedResponse
}
