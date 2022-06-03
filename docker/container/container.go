package container

import (
	"time"
)

// Container struct holds information about newly created container
type Container struct {
	ID        string        // ID of the container
	TimeLimit time.Duration // Maximum allowed time for the container to run
}
