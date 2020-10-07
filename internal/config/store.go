package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

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
		return "", err
	}

	configDirPath := filepath.Join(homeDir, ".gpair")
	if _, err = os.Stat(configDirPath); os.IsNotExist(err) {
		err := os.Mkdir(configDirPath, 0700)
		if err != nil {
			return "", err
		}
	}

	configPath := filepath.Join(configDirPath, "config.json")
	if _, err = os.Stat(configPath); os.IsNotExist(err) {
		internal.PrintVerbose("Creating new config at %s", configPath)
		file, err := os.Create(configPath)
		if err != nil {
			return "", err
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
		return false, err
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
		return err
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
		return err
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
		return NewConfig(), err
	}
	defer jsonFile.Close()

	jsonBytes, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return NewConfig(), err
	}

	if len(jsonBytes) == 0 {
		return NewConfig(), nil
	}

	var config Config
	// Since the config file is expected to be small, Marshal/Unmarshal should be adequate
	err = json.Unmarshal(jsonBytes, &config)
	if err != nil {
		return NewConfig(), err
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
	}

	return nil
}
