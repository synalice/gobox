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

	mount1, err := mount.NewMount(ctrl, "", "/userFolder1")
	if err != nil {
		t.Errorf("%v", err)
	}

	mount2, err := mount.NewMount(ctrl, "", "/userFolder2")
	if err != nil {
		t.Errorf("%v", err)
	}

	mount3, err := mount.NewMount(ctrl, "", "/userFolder3")
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
		CPUCount(6)
	containerConfig := configBuilder.Build()

	log.Println(containerConfig)
}
