package creation

import (
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	v1 "github.com/opencontainers/image-spec/specs-go/v1"
)

// Config contains data for creating a docker container
type Config struct {
	ContainerConfig  container.Config
	HostConfig       container.HostConfig
	NetworkingConfig network.NetworkingConfig
	Platform         v1.Platform
	TimeLimit        time.Duration
}
