package main

import (
	"fmt"
	"time"

	"github.com/synalice/gobox/docker/config"
	"github.com/synalice/gobox/docker/container"
	"github.com/synalice/gobox/docker/controller"
	"github.com/synalice/gobox/docker/volume"
)

func main() {
	c, err := controller.NewController()
	if err != nil {
		fmt.Printf("error creating new controller: %v\n", err)
	}

	volume1, err := volume.Ensure(c, "")
	if err != nil {
		fmt.Printf("error creating volume: %v\n", err)
	}
	defer func() {
		err = volume.Remove(c, volume1.Name)
		if err != nil {
			fmt.Printf("couldn't remove volume1: %v\n", err)
		}
	}()

	volume2, err := volume.Ensure(c, "")
	if err != nil {
		fmt.Printf("error creating volume: %v\n", err)
	}
	defer func() {
		err = volume.Remove(c, volume2.Name)
		if err != nil {
			fmt.Printf("couldn't remove volume2: %v\n", err)
		}
	}()

	volume3, err := volume.Ensure(c, "")
	if err != nil {
		fmt.Printf("error creating volume: %v\n", err)
	}
	defer func() {
		err = volume.Remove(c, volume3.Name)
		if err != nil {
			fmt.Printf("couldn't remove volume3: %v\n", err)
		}
	}()

	configBuilder := config.NewConfigBuilder(c)
	configBuilder.
		Image("python").
		Cmd("python").
		Mount(volume1, "/userFolder1").
		Mount(volume2, "/userFolder2").
		Mount(volume3, "/userFolder3").
		TimeLimit(1 * time.Second).
		MemoryLimit(64).
		CPUCount(1000).
		DiskSpace(1024)
	newConfig := configBuilder.Build()

	containerBuilder := container.NewContainerBuilder(c)
	containerBuilder.
		SetConfig(newConfig)
	builtContainer, err := containerBuilder.Build()
	if err != nil {
		fmt.Printf("couldn't build container: %v\n", err)
	}
	defer func() {
		err = container.Remove(c, builtContainer.ID)
		if err != nil {
			fmt.Printf("couldn't remove the container: %v\n", err)
		}
	}()

	err = container.Start(c, builtContainer.ID)
	if err != nil {
		fmt.Printf("couldn't start a container: %v\n", err)
	}

	_, err = container.Wait(c, builtContainer.ID, builtContainer.TimeLimit)
	if err != nil {
		if err.Error() == "context deadline exceeded" {
			fmt.Println("Container killed due to timeout")
			return
		} else {
			fmt.Printf("error while waiting for container to finish: %v\n", err)
		}
	}

	logs, err := container.GetLogs(c, builtContainer.ID)
	if err != nil {
		fmt.Printf("couldn't get container's logs: %v\n", err)
	}

	fmt.Println(logs)
}
