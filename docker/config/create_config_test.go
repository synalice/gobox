package config_test

import (
	"log"
	"testing"
	"time"

	"github.com/synalice/gobox/docker/config"
	"github.com/synalice/gobox/docker/controller"
	"github.com/synalice/gobox/docker/mount"
)

func TestConfigBuilder(t *testing.T) {
	ctrl, err := controller.NewController()
	if err != nil {
		t.Errorf("error creating new controller: %v", err)
	}

	mount1, err := mount.NewMount(ctrl, "", "/theFolder1")
	if err != nil {
		t.Errorf("%v", err)
	}

	mount2, err := mount.NewMount(ctrl, "", "/theFolder2")
	if err != nil {
		t.Errorf("%v", err)
	}

	mount3, err := mount.NewMount(ctrl, "", "/theFolder3")
	if err != nil {
		t.Errorf("%v", err)
	}

	configBuilder := config.NewConfigBuilder(ctrl)
	configBuilder.
		Image("python").
		Cmd("python").
		Mount(mount1).
		Mount(mount2).
		Mount(mount3).
		TimeLimit(3 * time.Second).
		MemoryLimit(64).
		CPUCount(6).
		DiskSpace(64)
	containerConfig := configBuilder.Build()

	err = mount.Remove(ctrl, mount1, mount2, mount3)
	if err != nil {
		t.Errorf("%v", err)
	}

	log.Println(containerConfig)
}
