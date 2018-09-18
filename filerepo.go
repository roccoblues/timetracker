package main

import (
	"io/ioutil"
	"os"

	"github.com/pkg/errors"
)

// fileRepo implements the repository interface with a simple file storage.
type fileRepo struct {
	path string
}

// newFileRepo returns a new fileRepo struct for the given path.
func newFileRepo(path string) *fileRepo {
	return &fileRepo{path: path}
}

// Read returns the contents of fileRepo.path. If the file at path doesn't
// exist error is nil and an empty byte slice is returned.
func (f *fileRepo) Read() ([]byte, error) {
	if _, err := os.Stat(f.path); err == nil {
		data, err := ioutil.ReadFile(f.path)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to read file '%s'", f.path)
		}
		return data, nil
	}

	return []byte{}, nil
}

// Write writes the given bytes to fileRepo.path creating the file
// if necessary.
func (f *fileRepo) Write(data []byte) error {
	if err := ioutil.WriteFile(f.path, data, 0644); err != nil {
		return errors.Wrapf(err, "failed to write to '%s'", f.path)
	}

	return nil
}
