package config

import (
	"fmt"
	"strings"
)

// ErrMissingCollaborator returns an error when an alias is requested that doesn't exist in the config
func ErrMissingCollaborator(missing []string) error {
	if len(missing) == 0 {
		return nil
	} else if len(missing) == 1 {
		return fmt.Errorf("No collaborator exists for the alias '%s'", missing[0])
	}

	return fmt.Errorf("No collaborators exist for aliases '%s'", strings.Join(missing, "', '"))
}

// ErrSaveFailure is returned when the config file could not be created
type ErrSaveFailure struct {
	Path string
	Err  error
}

func (err *ErrSaveFailure) Error() string {
	return fmt.Sprintf("failed to create config file at %s", err.Path)
}

func (err *ErrSaveFailure) Unwrap() error {
	return err.Err
}

// NewErrSaveFailure returns an instance of ErrSaveFailure
func NewErrSaveFailure(err error, path string) *ErrSaveFailure {
	return &ErrSaveFailure{Err: err, Path: path}
}
