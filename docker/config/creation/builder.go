// FIXME: Add config defaults!

package creation

import (
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/mount"
)

type Builder struct {
	Config *Config
}

// NewConfigBuilder return new builder that can be used to build config for a container
func NewConfigBuilder() *Builder {
	return &Builder{Config: &Config{}}
}

// Image sets image that container will use
func (b *Builder) Image(image string) *Builder {
	b.Config.ContainerConfig.Image = image
	return b
}

// Cmd sets command that will run when container starts
func (b *Builder) Cmd(cmd ...string) *Builder {
	b.Config.ContainerConfig.Cmd = cmd
	return b
}

// Mount sets a mount point that container will have. Can be called multiple times for multiple mounts
func (b *Builder) Mount(volume *types.Volume, containerPath string) *Builder {
	m := mount.Mount{
		Type:   mount.TypeVolume,
		Source: volume.Name,
		Target: containerPath,
	}
	b.Config.HostConfig.Mounts = append(b.Config.HostConfig.Mounts, m)

	return b
}

// TimeLimit sets max allowed time for the container to run
func (b *Builder) TimeLimit(timeLimit time.Duration) *Builder {
	b.Config.TimeLimit = timeLimit
	return b
}

// MemoryLimit sets max amount of memory allocated to the container (in megabytes)
func (b *Builder) MemoryLimit(memoryLimit int) *Builder {
	b.Config.HostConfig.Resources.Memory = int64(memoryLimit * 1024 * 1024)
	return b
}

// DiskSpace sets max disk space allocated for the container (in megabytes)
// TODO: This might not work. Not yet tested
func (b *Builder) DiskSpace(diskSpace int) *Builder {
	b.Config.HostConfig.StorageOpt = make(map[string]string)
	b.Config.HostConfig.StorageOpt["size"] = string(rune(diskSpace)) + "MB"
	return b
}

// CPUCount sets max amount of CPU cycles for the container.
// Be aware that it might collide with TimeLimit since 1 CPU cycle is approximately 1 second.
func (b *Builder) CPUCount(CPUCount int) *Builder {
	b.Config.HostConfig.CPUCount = int64(CPUCount)
	return b
}

// Build returns a config that describes a container
func (b *Builder) Build() *Config {
	return b.Config
}
