package docker_test

import (
	"testing"

	"github.com/synalice/gobox/docker"
)

var id string
var config = docker.ContainerConfig{
	Image:      "python",
	LocalImage: false,
	Cmd:        []string{"python"},
	Time:       3,
	MemoryMB:   64,
}

func TestNewController(t *testing.T) {
	_, err := docker.NewController(&config)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
}

func TestEnsureImage(t *testing.T) {
	c, err := docker.NewController(&config)
	if err != nil {
		t.Error(err)
	}

	err = c.EnsureImage(config.Image)
	if err != nil {
		t.Error(err)
	}
}

func TestContainerCreate(t *testing.T) {
	c, err := docker.NewController(&config)
	if err != nil {
		t.Error(err)
	}

	id, err = c.ContainerCreate(&config, []docker.VolumeMount{})
	if err != nil {
		t.Error(err)
	}
}

func TestContainerStart(t *testing.T) {
	c, err := docker.NewController(&config)
	if err != nil {
		t.Error(err)
	}

	err = c.ContainerStart(id)
	if err != nil {
		t.Error(err)
	}
}

func TestContainerWait(t *testing.T) {
	c, err := docker.NewController(&config)
	if err != nil {
		t.Error(err)
	}

	_, err = c.ContainerWait(id)
	if err != nil {
		t.Error(err)
	}
}

func TestContainerLog(t *testing.T) {
	c, err := docker.NewController(&config)
	if err != nil {
		t.Error(err)
	}

	data, err := c.ContainerLog(id)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(data)
}

func TestContainerRemove(t *testing.T) {
	c, err := docker.NewController(&config)
	if err != nil {
		t.Error(err)
	}

	err = c.ContainerRemove(id)
	if err != nil {
		t.Error(err)
	}
}

func TestRun(t *testing.T) {
	c, err := docker.NewController(&config)
	if err != nil {
		t.Error(err)
	}

	_, logs, err := c.Run([]docker.VolumeMount{})
	if err != nil {
		t.Error(err)
	}
	fmt.Println(logs)
}
