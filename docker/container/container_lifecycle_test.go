package container_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/synalice/gobox/docker/config"
	"github.com/synalice/gobox/docker/container"
	"github.com/synalice/gobox/docker/controller"
	"github.com/synalice/gobox/docker/volume"
)

func TestContainerLifecycle(t *testing.T) {
	c, err := controller.NewController()
	if err != nil {
		t.Errorf("error creating new controller: %v", err)
	}

	volume1, err := volume.Ensure(c, "")
	if err != nil {
		t.Errorf("error creating volume: %v", err)
	}

	volume2, err := volume.Ensure(c, "")
	if err != nil {
		t.Errorf("error creating volume: %v", err)
	}

	volume3, err := volume.Ensure(c, "")
	if err != nil {
		t.Errorf("error creating volume: %v", err)
	}

	configBuilder := config.NewConfigBuilder()
	configBuilder.
		Image("python").
		Cmd("python").
		Mount(volume1, "/userFolder1").
		Mount(volume2, "/userFolder2").
		Mount(volume3, "/userFolder3").
		TimeLimit(3 * time.Second).
		MemoryLimit(64).
		CPUCount(1000).
		DiskSpace(1024)
	newConfig := configBuilder.Build()

	containerBuilder := container.NewContainerBuilder(c)
	containerBuilder.
		SetConfig(newConfig)
	builtContainer, err := containerBuilder.Build()
	if err != nil {
		t.Errorf("couldn't build container: %v", err)
	}

	err = container.Start(c, builtContainer.ID)
	if err != nil {
		t.Errorf("couldn't start a container: %v", err)
	}

	_, err = container.Wait(c, builtContainer.ID, builtContainer.TimeLimit)
	if err != nil {
		if err.Error() == "context deadline exceeded" {
			fmt.Println("Container killed due to timeout")
		} else {
			t.Errorf("error while waiting for container to finish: %v", err)
		}
	}

	logs, err := container.GetLogs(c, builtContainer.ID)
	if err != nil {
		t.Errorf("couldn't get container's logs: %v", err)
	}

	err = container.Remove(c, builtContainer.ID)
	if err != nil {
		t.Errorf("couldn't remove the container: %v", err)
	}

	err = volume.Remove(c, volume1.Name)
	if err != nil {
		t.Errorf("couldn't remove volume1: %v", err)
	}

	err = volume.Remove(c, volume2.Name)
	if err != nil {
		t.Errorf("couldn't remove volume2: %v", err)
	}

	err = volume.Remove(c, volume3.Name)
	if err != nil {
		t.Errorf("couldn't remove volume3: %v", err)
	}

	fmt.Println(logs)
}
