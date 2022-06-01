package container_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/synalice/gobox/docker"
	"github.com/synalice/gobox/docker/config"
	"github.com/synalice/gobox/docker/container"
)

func TestBuilder(t *testing.T) {
	c, err := docker.NewController()
	if err != nil {
		t.Errorf("error creating new controller: %v", err)
	}

	volume1, err := c.EnsureVolume("")
	if err != nil {
		t.Errorf("error creating volume: %v", err)
	}
	defer func(c *docker.Controller, name string) {
		err := c.RemoveVolume(name)
		if err != nil {
			t.Errorf("couldn't remove volume: %v", err)
		}
	}(c, volume1.Name)
	if err != nil {
		t.Errorf("couldn't remove volume: %v", err)
	}

	volume2, err := c.EnsureVolume("")
	if err != nil {
		t.Errorf("error creating volume: %v", err)
	}
	defer func(c *docker.Controller, name string) {
		err := c.RemoveVolume(name)
		if err != nil {
			t.Errorf("couldn't remove volume: %v", err)
		}
	}(c, volume2.Name)
	if err != nil {
		t.Errorf("couldn't remove volume: %v", err)
	}

	volume3, err := c.EnsureVolume("")
	if err != nil {
		t.Errorf("error creating volume: %v", err)
	}
	defer func(c *docker.Controller, name string) {
		err := c.RemoveVolume(name)
		if err != nil {
			t.Errorf("couldn't remove volume: %v", err)
		}
	}(c, volume3.Name)
	if err != nil {
		t.Errorf("couldn't remove volume: %v", err)
	}

	configBuilder := config.NewConfigBuilder()
	configBuilder.
		Image("python").
		Cmd("python", "...").
		Mount(volume1, "/userFolder1").
		Mount(volume2, "/userFolder2").
		Mount(volume3, "/userFolder3").
		TimeLimit(3 * time.Second).
		MemoryLimit(64).
		CPUCount(6).
		DiskSpace(1024)
	newConfig := configBuilder.Build()

	containerBuilder := container.NewContainerBuilder(c)
	containerBuilder.
		SetConfig(newConfig)
	builtContainer, err := containerBuilder.Build()
	if err != nil {
		t.Errorf("couldn't build container: %v", err)
	}

	fmt.Print("Container ID:", builtContainer.ID)
}
