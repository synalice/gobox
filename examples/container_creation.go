package main

import (
	"fmt"
	"log"
	"time"

	"github.com/synalice/gobox/docker/config"
	"github.com/synalice/gobox/docker/container"
	"github.com/synalice/gobox/docker/controller"
	"github.com/synalice/gobox/docker/file"
	"github.com/synalice/gobox/docker/mount"
)

func main() {
	ctrl, err := controller.NewController()
	if err != nil {
		log.Fatalln(err)
	}

	mount1, err := mount.NewMount(ctrl, "", "/userFolder1")
	if err != nil {
		log.Fatalln(err)
	}

	mount2, err := mount.NewMount(ctrl, "", "/userFolder2")
	if err != nil {
		log.Fatalln(err)
	}

	mount3, err := mount.NewMount(ctrl, "", "/userFolder3")
	if err != nil {
		log.Fatalln(err)
	}

	myFile := file.File{
		Name: "main.py",
		Body: "print(\"Hello, World!\")",
	}

	configBuilder := config.NewConfigBuilder(ctrl)
	configBuilder.
		Image("python").
		Cmd("python", "/userFolder1/main.py").
		Mount(mount1).
		Mount(mount2).
		Mount(mount3).
		TimeLimit(1 * time.Second).
		MemoryLimit(64).
		CPUCount(1000).
		DiskSpace(1024)
	newConfig := configBuilder.Build()

	containerBuilder := container.NewContainerBuilder(ctrl)
	containerBuilder.
		SetConfig(newConfig).
		SetFile(myFile, mount1)
	builtContainer, err := containerBuilder.Build()
	if err != nil {
		log.Fatalln(err)
	}

	err = container.Start(ctrl, builtContainer.ID)
	if err != nil {
		log.Fatalln(err)
	}

	_, err = container.Wait(ctrl, builtContainer.ID, builtContainer.TimeLimit)
	if err.Error() == "container killed due to timeout" {
		log.Println(err)
	} else {
		log.Fatalln(err)
	}

	logs, err := container.GetLogs(ctrl, builtContainer.ID)
	if err != nil {
		log.Fatalln(err)
	}

	err = container.Remove(ctrl, builtContainer.ID)
	if err != nil {
		log.Fatalln(err)
	}

	err = mount.CleanUp(ctrl, mount1, mount2, mount3)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(logs)
}
