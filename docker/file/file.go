package file

// File defines a file that will be executed inside the container
type File struct {
	Name string // Name of the file
	Body string // Body of the file. The code that will be executed
}
