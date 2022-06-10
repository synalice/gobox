package container_test

import (
	"fmt"
	"log"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/synalice/gobox/docker/config"
	"github.com/synalice/gobox/docker/container"
	"github.com/synalice/gobox/docker/controller"
	"github.com/synalice/gobox/docker/file"
	"github.com/synalice/gobox/docker/mount"
)

var wg sync.WaitGroup

func TestOnLoad(t *testing.T) {
	ctrl, err := controller.NewController()
	if err != nil {
		t.Errorf("%v", err)
	}

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(iter int) {
			log.Println("Starting goroutine " + strconv.Itoa(iter))
			mount1, err := mount.NewMount(ctrl, "", "/userFolder1")
			if err != nil {
				t.Errorf("%v", err)
			}

			myFile := file.File{
				Name: "main.py",
				Body: "print(\"Hello, Goroutine " + strconv.Itoa(iter) + "!\")",
			}

			configBuilder := config.NewConfigBuilder(ctrl)
			configBuilder.
				Image("python").
				Cmd("python", "/userFolder1/main.py").
				Mount(mount1).
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
				t.Errorf("%v", err)
			}

			log.Println("Starting container " + strconv.Itoa(iter))
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

			fmt.Println(logs)

			err = mount.Remove(ctrl, mount1)
			if err != nil {
				t.Errorf("%v", err)
			}

			wg.Done()
		}(i)
	}

	wg.Wait()
}
