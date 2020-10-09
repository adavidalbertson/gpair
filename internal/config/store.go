package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/pkg/errors"

	"github.com/adavidalbertson/gpair/internal"
)

// store represents some persistent storage for config
type store interface {
	load() (Config, error)
	save(Config) error
}

// fileStore is an implementation of store that reads from and writes to a file (~/.gpair/config.json)
type fileStore struct {
	path string
}

func (fs *fileStore) configPath() (string, error) {
	if fs.path != "" {
		return fs.path, nil
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", errors.Wrap(err, "failed to locate home directory")
	}

	configDirPath := filepath.Join(homeDir, ".gpair")
	if _, err = os.Stat(configDirPath); os.IsNotExist(err) {
		internal.PrintVerbose("Creating new config directory at %s", configDirPath)
		err := os.Mkdir(configDirPath, 0700)
		if err != nil {
			return "", NewErrSaveFailure(err, configDirPath)
		}
	}

	configPath := filepath.Join(configDirPath, "config.json")
	if _, err = os.Stat(configPath); os.IsNotExist(err) {
		internal.PrintVerbose("Creating new config file at %s", configPath)
		file, err := os.Create(configPath)
		if err != nil {
			return "", NewErrSaveFailure(err, configPath)
		}
		defer file.Close()
	}

	fs.path = configPath
	return fs.path, nil
}

func (fs *fileStore) fileExists() (bool, error) {
	path, err := fs.configPath()
	if err != nil {
		return false, err
	}

	_, err = os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		
		if os.IsPermission(err) {
			return false, NewErrSaveFailure(err, path)
		}

		return false, errors.Wrapf(err, "failed to determine existence of config file at %s", path)
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

	path, err := fs.configPath()
	if err != nil {
		return err
	}

	err = os.Rename(path, path+".bak")
	if err != nil {
		return errors.Wrap(err, "failed to create config file backup")
	}

	return nil
}

func (fs *fileStore) restoreBackup() error {
	path, err := fs.configPath()
	if err != nil {
		return err
	}

	err = os.Rename(path+".bak", path)
	if err != nil {
		return errors.Wrap(err, "failed to restore config file backup")
	}

	return nil
}

func (fs *fileStore) load() (Config, error) {
	path, err := fs.configPath()
	if err != nil {
		return NewConfig(), err
	}

	jsonFile, err := os.Open(path)
	if err != nil {
		err = fs.createBackup()
		if err != nil {
			return NewConfig(), err
		}

		internal.PrintVerbose("Failed to open config file at %s. Starting a new config file. Original has been moved to %s.bak", path, path)
		return NewConfig(), nil
	}
	defer jsonFile.Close()

	jsonBytes, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		err = fs.createBackup()
		if err != nil {
			return NewConfig(), err
		}

		internal.PrintVerbose("Failed to parse config file at %s. Starting a new config file. Original has been moved to %s.bak", path, path)
		return NewConfig(), nil
	}

	if len(jsonBytes) == 0 {
		return NewConfig(), nil
	}

	var config Config
	// Since the config file is expected to be small, Marshal/Unmarshal should be adequate
	err = json.Unmarshal(jsonBytes, &config)
	if err != nil {
		err = fs.createBackup()
		if err != nil {
			return NewConfig(), err
		}

		internal.PrintVerbose("Failed to parse config file at %s. Starting a new config file. Original has been moved to %s.bak", path, path)
		return NewConfig(), nil
	}

	return config, nil
}

func (fs *fileStore) save(config Config) error {
	jsonBytes, err := json.Marshal(config)
	if err != nil {
		return err
	}

	err = fs.createBackup()
	if err != nil {
		return err
	}

	path, err := fs.configPath()
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(path, jsonBytes, 0700)
	if err != nil {
		err = fs.restoreBackup()
		if err != nil {
			return err
		}
		internal.PrintVerbose("Failed to save config, backup has been restored.")
	}

	return nil
}
