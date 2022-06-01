package config_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/synalice/gobox/docker"
	creationconfig "github.com/synalice/gobox/docker/config/creation"
)

func TestConfigBuilder(t *testing.T) {
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

	configBuilder := creationconfig.NewConfigBuilder()
	configBuilder.
		Image("python").
		Cmd("python", "...").
		Mount(volume1, "/userFolder").
		Mount(volume2, "/userFolder").
		Mount(volume3, "/userFolder").
		TimeLimit(3 * time.Second).
		MemoryLimit(64).
		CPUCount(6).
		DiskSpace(64)
	containerConfig := configBuilder.Build()

	fmt.Println(containerConfig)
}
