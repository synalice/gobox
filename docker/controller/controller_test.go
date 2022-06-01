package controller_test

import (
	"testing"

	"github.com/synalice/gobox/docker/controller"
)

func TestNewController(t *testing.T) {
	_, err := controller.NewController()
	if err != nil {
		t.Error("couldn't create new Controller object")
	}
}
