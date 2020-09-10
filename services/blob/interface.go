package blob

import "io"

// Large (more then memory) content storage.
type Blob interface {
	// Push content to the storage using provided writer.
	Push(id string, handler func(out io.Writer) error) error
	// Get content from the storage.
	Get(id string) (io.ReadCloser, error)
}
