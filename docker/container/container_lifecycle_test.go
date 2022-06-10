package container_test

import (
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/synalice/gobox/docker/config"
	"github.com/synalice/gobox/docker/container"
	"github.com/synalice/gobox/docker/controller"
	"github.com/synalice/gobox/docker/file"
	"github.com/synalice/gobox/docker/mount"
)

func TestContainerLifecycle(t *testing.T) {
	ctrl, err := controller.NewController()
	if err != nil {
		t.Errorf("%v", err)
	}

	myFile := file.File{
		Name: "main.py",
		Body: "print(\"Hello, World!\")",
	}

	mount1, err := mount.NewMount(ctrl, "", "/theFolder1")
	if err != nil {
		t.Errorf("%v", err)
	}

	mount2, err := mount.NewMount(ctrl, "", "/theFolder2")
	if err != nil {
		t.Errorf("%v", err)
	}

	configBuilder := config.NewConfigBuilder(ctrl)
	configBuilder.
		Image("python").
		Cmd("python", "/userFolder1/main.py").
		Mount(mount1).
		Mount(mount2).
		TimeLimit(1 * time.Second).
		MemoryLimit(64)
	newConfig := configBuilder.Build()

	containerBuilder := container.NewContainerBuilder(ctrl)
	containerBuilder.
		SetConfig(newConfig).
		SetFile(myFile, mount1)
	builtContainer, err := containerBuilder.Build()
	if err != nil {
		t.Errorf("%v", err)
	}

	err = container.Start(ctrl, builtContainer.ID)
	if err != nil {
		t.Errorf("%v", err)
	}

	_, err = container.Wait(ctrl, builtContainer.ID, builtContainer.TimeLimit)
	if err == container.ErrorTimeout {
		log.Println(err)
	} else if err != nil {
		t.Errorf("%v", err)
	}

	logs, err := container.GetLogs(ctrl, builtContainer.ID)
	if err != nil {
		t.Errorf("%v", err)
	}

	err = container.Remove(ctrl, builtContainer.ID)
	if err != nil {
		t.Errorf("%v", err)
	}

	err = mount.CleanUp(ctrl, mount1, mount2)
	if err != nil {
		t.Errorf("%v", err)
	}

	fmt.Println(logs)
}
