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
	// First, we need to create a new controller object. It is needed
	// everywhere, where we work with a docker's API.
	ctrl, err := controller.NewController()
	if err != nil {
		log.Println(err)
	}

	// Let's describe the file that we want to execute inside our container.
	// It needs a name and a body - the actual code, that we want to
	// run.
	myFile := file.File{
		Name: "main.py",
		Body: "print(\"Hello, World!\")",
	}

	// Now, let's create a mount. It specifies where do we want to create
	// our file.
	//
	// It is done by creating a volume that connects to a specific
	// folder inside the container and then our file is created in this
	// folder. This ensures that each container has at least 1 volume
	// associated with it.
	mount1, err := mount.NewMount(ctrl, "", "/theFolder1")
	if err != nil {
		log.Println(err)
	}

	// We can create any amount of mounts so let's create one more.
	mount2, err := mount.NewMount(ctrl, "", "/theFolder2")
	if err != nil {
		log.Println(err)
	}

	// config.NewConfigBuilder() is an implementation of a builder pattern.
	// It is responsible for creating a config that describes how should
	// a future container look like.
	configBuilder := config.NewConfigBuilder(ctrl)
	configBuilder.
		// An image that the container will use.
		Image("python").

		// A command that will be executed when container starts.
		Cmd("python", "/theFolder1/main.py").

		// Mounting the mounts! We can set how many we want.
		Mount(mount1).
		Mount(mount2).

		// Maximum amount of time for the container to run.
		TimeLimit(1 * time.Second).

		// How much RAM (in megabytes [mebibytes]) is allocated for the
		// container.
		MemoryLimit(64)

	// Now we let's build a config.
	newConfig := configBuilder.Build()

	// config.NewConfigBuilder() is also an implementation of a builder
	// pattern. It is used to build an actual container.
	containerBuilder := container.NewContainerBuilder(ctrl)
	containerBuilder.
		// A config that the container should have
		SetConfig(newConfig).

		// And a file that will be created.
		//
		// Notice how you need to specify a file and then a mount. The
		// file will be created in a folder that this mount has - in
		// this situation it is /theFolder1
		SetFile(myFile, mount1)

	// And now we actually create a container. It will appear on your
	// system.
	builtContainer, err := containerBuilder.Build()
	if err != nil {
		log.Println(err)
	}

	// Let's start this container.
	err = container.Start(ctrl, builtContainer.ID)
	if err != nil {
		log.Println(err)
	}

	// Let's wait until the container finishes its work.
	//
	// One of the errors you might actually want to expect from this
	// function is container.ErrorTimeout. If you get it - it means that
	// the container hasn't finished it's work in the allotted time and was
	// killed.
	_, err = container.Wait(ctrl, builtContainer.ID, builtContainer.TimeLimit)
	if err == container.ErrorTimeout {
		log.Println(err)
	} else if err != nil {
		log.Println(err)
	}

	// Let's get the logs from the container's stdout.
	logs, err := container.GetLogs(ctrl, builtContainer.ID)
	if err != nil {
		log.Println(err)
	}

	// Remove the container.
	err = container.Remove(ctrl, builtContainer.ID)
	if err != nil {
		log.Println(err)
	}

	// Remove volumes
	err = mount.Remove(ctrl, mount1, mount2)
	if err != nil {
		log.Println(err)
	}

	// And finally - let's print the results of container's work.
	fmt.Println(logs)
}
