package config

import (
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	v1 "github.com/opencontainers/image-spec/specs-go/v1"
)

// Config contains fields necessary for creating a docker container
type Config struct {
	//Fields related to docker API
	ContainerConfig  container.Config
	HostConfig       container.HostConfig
	NetworkingConfig network.NetworkingConfig
	Platform         v1.Platform

	// Maximum allowed time for the container to run
	TimeLimit time.Duration
}
