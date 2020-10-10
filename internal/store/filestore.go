package store

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/pkg/errors"

	"github.com/adavidalbertson/gpair/internal"
)

// Options for where to start the path
const (
	HOME = iota
	CACHE
	CONFIG
	TEMP
	ROOT
)

// Store represents an object that can store bytes, typically a file on disk
type Store interface {
	Read() ([]byte, error)
	Write([]byte) error
	GetPath() string
}

// fileStore is an implementation of store that reads from and writes to a file (~/.gpair/config.json)
type fileStore struct {
	path string
}

// NewFileStore returns a filestore backed by a file, which it creates if it does not exist already
func NewFileStore(filename string, startDirType int, dirPath ...string) (Store, error) {
	fs := &fileStore{}

	path, err := getStartDir(startDirType)
	if err != nil {
		return nil, err
	}

	for _, dir := range dirPath {
		path = filepath.Join(path, dir)
		if _, err = os.Stat(path); os.IsNotExist(err) {
			internal.PrintVerbose("Creating new directory at %s", path)
			err := os.Mkdir(path, 0700)
			if err != nil {
				return nil, NewErrFileInaccessible(err, path)
			}
		}
	}

	path = filepath.Join(path, filename)
	if _, err = os.Stat(path); os.IsNotExist(err) {
		internal.PrintVerbose("Creating new file at %s", path)
		file, err := os.Create(path)
		if err != nil {
			return nil, NewErrFileInaccessible(err, path)
		}
		defer file.Close()
	}

	fs.path = path
	return fs, nil
}

func getStartDir(startDirType int) (startDir string, err error) {
	switch startDirType {
	case HOME:
		startDir, err = os.UserHomeDir()
		if err != nil {
			return "", errors.Wrap(err, "failed to locate user home directory")
		}

	case CACHE:
		startDir, err = os.UserCacheDir()
		if err != nil {
			return "", errors.Wrap(err, "failed to locate user cache directory")
		}

	case CONFIG:
		startDir, err = os.UserConfigDir()
		if err != nil {
			return "", errors.Wrap(err, "failed to locate user config directory")
		}

	case TEMP:
		return os.TempDir(), nil

	case ROOT:
		return "/", nil

	default:
		return "", errors.New("not a recognized start directory type")
	}

	return
}

func (fs *fileStore) fileExists() (bool, error) {
	_, err := os.Stat(fs.path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}

		if os.IsPermission(err) {
			return false, NewErrFileInaccessible(err, fs.path)
		}

		return false, errors.Wrapf(err, "failed to determine existence of file at %s", fs.path)
	}

	return true, nil
}

func (fs *fileStore) createBackup() error {
	exists, err := fs.fileExists()
	if err != nil {
		return err
	}

	if !exists {
		return nil
	}

	err = os.Rename(fs.path, fs.path+".bak")
	if err != nil {
		return errors.Wrapf(err, "failed to create backup file %s.bak", fs.path)
	}

	return nil
}

func (fs *fileStore) restoreBackup() error {
	err := os.Rename(fs.path+".bak", fs.path)
	if err != nil {
		return errors.Wrapf(err, "failed to restore backup file at %s.bak", fs.path)
	}

	return nil
}

func (fs *fileStore) Read() ([]byte, error) {
	jsonFile, err := os.Open(fs.path)
	if err != nil {
		err = fs.createBackup()
		if err != nil {
			return nil, err
		}

		internal.PrintVerbose("Failed to open file at %s. Starting an empty file. Original has been moved to %s.bak", fs.path, fs.path)
		return nil, nil
	}
	defer jsonFile.Close()

	jsonBytes, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		err = fs.createBackup()
		if err != nil {
			return nil, err
		}

		internal.PrintVerbose("Failed to parse file at %s. Starting an empty file. Original has been moved to %s.bak", fs.path, fs.path)
		return nil, nil
	}

	return jsonBytes, nil
}

func (fs *fileStore) Write(jsonBytes []byte) error {
	err := fs.createBackup()
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(fs.path, jsonBytes, 0700)
	if err != nil {
		err = fs.restoreBackup()
		if err != nil {
			return err
		}
		internal.PrintVerbose("Failed to save file at %s, backup has been restored.", fs.path)
	}

	return nil
}

func (fs *fileStore) GetPath() string {
	return fs.path
}
