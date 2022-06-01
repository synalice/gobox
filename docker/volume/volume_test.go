package volume_test

import (
	"testing"

	"github.com/synalice/gobox/docker/controller"
	"github.com/synalice/gobox/docker/volume"
)

func TestVolumeLifecycle(t *testing.T) {
	c, err := controller.NewController()
	if err != nil {
		t.Error(err)
	}

	volume1, err := volume.Ensure(c, "myVolume")
	if err != nil {
		t.Errorf("couldn't create a volume")
	}

	volume2, err := volume.Ensure(c, "myVolume")
	if volume2.Name != volume1.Name {
		t.Errorf("should have used an already existing volume")
	}

	err = volume.Remove(c, "myVolume")
	if err != nil {
		t.Errorf("couldn't remove a volume")
	}
}
