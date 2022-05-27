package container_creation

import (
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	v1 "github.com/opencontainers/image-spec/specs-go/v1"
)

type Container struct {
	hostConfig       container.HostConfig
	networkingConfig network.NetworkingConfig
	platform         v1.Platform
	mounts           []mount.Mount
}
