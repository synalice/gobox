package file_test

import (
	"testing"

	"github.com/synalice/gobox/docker/file"
)

func TestCreateTarball(t *testing.T) {
	files := []file.File{
		{Name: "file_one.txt", Body: "this is a text-file"},
		{Name: "file_two.py", Body: "print(\"this a a python file\")"},
	}
	_, err := file.CreateTarball(files)

	if err != nil {
		t.Errorf("%v", err)
	}
}
