package docker

import (
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

// ContainerConfig describes what configuration the container should have
type ContainerConfig struct {
	ForBuild   bool          // Determines if the config is meant to be used
	Image      string        // Image which will be used for running code
	LocalImage bool          // Set to true if image shouldn't be pulled from the outside
	Cmd        []string      // Defines what will be run when container starts. Use "..." as a name of the file
	TimeLimit  time.Duration // Max time for the container to run
	MemoryMB   int64         // Container memory limit (in megabytes)
}

// ContainerFile describes what file the container should execute
type ContainerFile struct {
	FileName string // Name of file that will contain all the Content
	Content  string // Code that will be executed
}

// Controller is an object that wll be used for running methods off of it
type Controller struct {
	Cli *client.Client
}

// VolumeMount is used for specifying volumes for the container to mount
type VolumeMount struct {
	HostPath string
	Volume   *types.Volume
}
