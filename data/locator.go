package data

import (
	"io"
	"os"
)

// Locator describes the place where the bytes of the media are actually
// stored.  Locator is an interface so different saved data can be returned
// for the same metadata.  For example, there may be duplicates of an
// image stored at different locations.
type Locator interface {
	// Source indicates the type of location, e.g. "file" is a file on the
	// filesystem.
	Source() string
	// Data opens a stream to retrieve the contents of the image at that
	// location.
	Data() (io.ReadCloser, error)
}

type fileSystemLocator struct {
	Path string
}

func (fsl fileSystemLocator) Source() string {
	return "files"
}

func (fsl fileSystemLocator) Data() (io.ReadCloser, error) {
	return os.Open(fsl.Path)
}
