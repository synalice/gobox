package container_creation

import (
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	v1 "github.com/opencontainers/image-spec/specs-go/v1"
)

type ContainerConfigBuilder struct {
	ContainerBuilder
}

func (b *ContainerConfigBuilder) SetHostConfig(hc container.HostConfig) *ContainerConfigBuilder {
	b.container.hostConfig = hc
	return b
}

func (b *ContainerConfigBuilder) SetNetworkingConfig(nc network.NetworkingConfig) *ContainerConfigBuilder {
	b.container.networkingConfig = nc
	return b
}

func (b *ContainerConfigBuilder) SetPlatform(pl v1.Platform) *ContainerConfigBuilder {
	b.container.platform = pl
	return b
}
