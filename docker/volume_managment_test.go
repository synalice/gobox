package docker_test

import (
	"testing"

	"github.com/synalice/gobox/docker"
)

func TestVolumeLifecycle(t *testing.T) {
	c, err := docker.NewController()
	if err != nil {
		t.Error(err)
	}

	created, _, err := c.EnsureVolume("myVolume")
	if created != true {
		t.Errorf("Should have created the volume the first time")
	}

	created, _, err = c.EnsureVolume("myVolume")
	if created != false {
		t.Errorf("Should not have created the volume the second time")
	}

	removed, err := c.RemoveVolume("myVolume")
	if removed != true {
		t.Errorf("Should have removed the volume")
	}
}
