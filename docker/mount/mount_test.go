package mount_test

import (
	"testing"

	"github.com/synalice/gobox/docker/controller"
	"github.com/synalice/gobox/docker/mount"
)

func TestNewMount(t *testing.T) {
	ctrl, err := controller.NewController()
	if err != nil {
		t.Errorf("%v", err)
	}

	myMount, err := mount.NewMount(ctrl, "some_volume", "/someFolder")
	if err != nil {
		t.Errorf("%v", err)
	}

	err = mount.CleanUp(ctrl, myMount)
	if err != nil {
		t.Errorf("%v", err)
	}
}
