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
		log.Println(err)
	}

	myFile := file.File{
		Name: "main.py",
		Body: `
n = int(input())

if n % 3 == 0 and n % 5 == 0:
    print("FizzBuzz")
elif n % 3 == 0:
    print("Fizz")
elif n % 5 == 0:
    print("Buzz")
else:
    print(n)
`,
	}

	mount1, err := mount.NewMount(ctrl, "", "/theFolder1")
	if err != nil {
		log.Println(err)
	}

	configBuilder := config.NewConfigBuilder(ctrl)

	configBuilder.
		Image("python").
		Cmd("python", "/theFolder1/main.py").
		Mount(mount1).
		TimeLimit(1 * time.Second).
		MemoryLimit(64)

	newConfig := configBuilder.Build()

	containerBuilder := container.NewContainerBuilder(ctrl)

	containerBuilder.
		SetConfig(newConfig).
		SetFile(myFile, mount1)

	builtContainer, err := containerBuilder.Build()
	if err != nil {
		log.Println(err)
	}

	err = container.Start(ctrl, builtContainer, "15")
	if err != nil {
		log.Println(err)
	}

	_, err = container.Wait(ctrl, builtContainer.ID, builtContainer.TimeLimit)
	if err != nil {
		log.Println(err)
	}

	logs, err := container.GetLogs(ctrl, builtContainer)
	if err != nil {
		log.Println(err)
	}

	err = container.Remove(ctrl, builtContainer.ID)
	if err != nil {
		log.Println(err)
	}

	err = mount.Remove(ctrl, mount1)
	if err != nil {
		log.Println(err)
	}

	fmt.Println(logs)
}
