package fsblob

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"nano-run/services/blob"
)

// Dummy file storage: large storage based on file system.
// ID - just name of file, but content will be written atomically.
func New(rootDir string) blob.Blob {
	return &fsBlob{
		rootDir: rootDir,
		checker: func(id string) bool {
			return true
		},
	}
}

func NewCheck(rootDir string, checkFn func(string) bool) blob.Blob {
	return &fsBlob{
		rootDir: rootDir,
		checker: checkFn,
	}
}

type fsBlob struct {
	rootDir string
	checker func(id string) bool
}

func (k *fsBlob) Push(id string, handler func(out io.Writer) error) error {
	if !k.checker(id) {
		return fmt.Errorf("push: invalid id %s", id)
	}
	dir := filepath.Join(k.rootDir, id[0:1])
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return err
	}
	tempFile := filepath.Join(dir, id+".temp")
	destFile := filepath.Join(dir, id)

	tempF, err := os.Create(tempFile)
	if err != nil {
		return fmt.Errorf("fsblob: put data: create temp file: %w", err)
	}
	err = handler(tempF)
	flushErr := tempF.Close()
	if err != nil {
		_ = os.Remove(tempFile)
		return err
	}
	if flushErr != nil {
		_ = os.Remove(tempFile)
		return fmt.Errorf("fsblob: put data: flush content to temp file: %w", err)
	}

	err = os.Rename(tempFile, destFile)
	if err != nil {
		return fmt.Errorf("fsblob: put data: commit file: %w", err)
	}
	return nil
}

func (k *fsBlob) Get(id string) (io.ReadCloser, error) {
	if !k.checker(id) {
		return nil, fmt.Errorf("get: invalid id %s", id)
	}
	dir := filepath.Join(k.rootDir, id[0:1])
	destFile := filepath.Join(dir, id)
	f, err := os.Open(destFile)
	if err == nil {
		return f, nil
	}
	return nil, fmt.Errorf("fsblob: get file: %w", err)
}
