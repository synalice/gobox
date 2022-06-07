package file

import (
	"archive/tar"
	"bytes"
	"fmt"
)

// CreateTarball returns an in-memory tar archive of files
func CreateTarball(files []File) (bytes.Buffer, error) {
	var buf bytes.Buffer

	tw := tar.NewWriter(&buf)

	for _, f := range files {
		hdr := &tar.Header{
			Name: f.Name,
			Mode: 0700,
			Size: int64(len(f.Body)),
		}

		err := tw.WriteHeader(hdr)
		if err != nil {
			return buf, fmt.Errorf("error creating a tarball: %w", err)
		}

		_, err = tw.Write([]byte(f.Body))
		if err != nil {
			return buf, fmt.Errorf("error creating a tarball: %w", err)
		}
	}

	err := tw.Close()
	if err != nil {
		return buf, fmt.Errorf("error creating a tarball: %w", err)
	}

	return buf, nil
}
