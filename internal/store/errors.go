package store

import "fmt"

// ErrFileInaccessible is returned when a file could not be found or created
type ErrFileInaccessible struct {
	Path string
	Err  error
}

func (err *ErrFileInaccessible) Error() string {
	return fmt.Sprintf("failed to create config file at %s", err.Path)
}

func (err *ErrFileInaccessible) Unwrap() error {
	return err.Err
}

// NewErrFileInaccessible returns an instance of ErrFileInaccessible
func NewErrFileInaccessible(err error, path string) *ErrFileInaccessible {
	return &ErrFileInaccessible{Err: err, Path: path}
}