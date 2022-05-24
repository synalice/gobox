package docker

import (
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

// ContainerConfig describes what configuration the container should have
type ContainerConfig struct {
	Image      string        // Image which will be used for running code
	LocalImage bool          // Set to true if image shouldn't be pulled from the outside
	Cmd        []string      // Defines what will be run when container starts. Use "..." as name of the file
	Time       time.Duration // Max time for the code to execute FIXME: Unused in code
	MemoryMB   int64         // Max amount of memory for the container FIXME: Unused in code
}

// ContainerFile describes what file the container should execute
type ContainerFile struct {
	FileName string // Name of file that will contain all the Content
	Content  string // Code that will be executed
}

// Controller is an object that wll be used for running methods off of it
type controller struct {
	cli        *client.Client
	image      string
	localImage bool
	cmd        []string
}

type limits struct {
	time     time.Duration
	memoryMB int64
}

// VolumeMount TODO: Abstract volume creation and mounting away from user. This struct should be hidden
type VolumeMount struct {
	HostPath string
	Volume   *types.Volume
}
