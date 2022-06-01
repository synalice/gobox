package config_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/synalice/gobox/docker"
	"github.com/synalice/gobox/docker/config"
	"github.com/synalice/gobox/docker/controller"
	"github.com/synalice/gobox/docker/volume"
)

func TestConfigBuilder(t *testing.T) {
	c, err := docker.NewController()
	if err != nil {
		t.Errorf("error creating new controller: %v", err)
	}

	volume1, err := volume.EnsureVolume(c, "")
	if err != nil {
		t.Errorf("error creating volume: %v", err)
	}
	defer func(c *controller.Controller, name string) {
		err := volume.RemoveVolume(c, name)
		if err != nil {
			t.Errorf("couldn't remove volume1: %v", err)
		}
	}(c, volume1.Name)

	volume2, err := volume.EnsureVolume(c, "")
	if err != nil {
		t.Errorf("error creating volume: %v", err)
	}
	defer func(c *controller.Controller, name string) {
		err := volume.RemoveVolume(c, name)
		if err != nil {
			t.Errorf("couldn't remove volume2: %v", err)
		}
	}(c, volume2.Name)

	volume3, err := volume.EnsureVolume(c, "")
	if err != nil {
		t.Errorf("error creating volume: %v", err)
	}
	defer func(c *controller.Controller, name string) {
		err := volume.RemoveVolume(c, name)
		if err != nil {
			t.Errorf("couldn't remove volume3: %v", err)
		}
	}(c, volume3.Name)

	configBuilder := config.NewConfigBuilder()
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
