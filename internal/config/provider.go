package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/adavidalbertson/gpair/internal"
)

type provider interface {
	load() (config, error)
	save(config) error
}

type fileProvider struct {
	absoluteConfigPath string
}

func (fp *fileProvider) configPath() (string, error) {
	if fp.absoluteConfigPath != "" {
		return fp.absoluteConfigPath, nil
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	configDirPath, err := filepath.Abs(filepath.Join(homeDir, ".gpair"))
	if err != nil {
		return "", err
	}

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

	fp.absoluteConfigPath = configPath
	return fp.absoluteConfigPath, nil
}

func (fp *fileProvider) fileExists() (bool, error) {
	path, err := fp.configPath()
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

func (fp *fileProvider) createBackup() error {
	exists, err := fp.fileExists()
	if err != nil {
		return err
	}

	if !exists {
		return nil
	}

	path, err := fp.configPath()
	if err != nil {
		return err
	}

	err = os.Rename(path, path+".bak")
	if err != nil {
		return err
	}

	return nil
}

func (fp *fileProvider) restoreBackup() error {
	path, err := fp.configPath()
	if err != nil {
		return err
	}

	err = os.Rename(path+".bak", path)
	if err != nil {
		return err
	}

	return nil
}

func (fp *fileProvider) load() (config, error) {
	path, err := fp.configPath()
	if err != nil {
		return newConfig(), err
	}

	jsonFile, err := os.Open(path)
	if err != nil {
		return newConfig(), err
	}
	defer jsonFile.Close()

	jsonBytes, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return newConfig(), err
	}

	if len(jsonBytes) == 0 {
		return newConfig(), nil
	}

	var config config
	// Since the config file is expected to be small, Marshal/Unmarshal should be adequate
	err = json.Unmarshal(jsonBytes, &config)
	if err != nil {
		return newConfig(), err
	}

	return config, nil
}

func (fp *fileProvider) save(config config) error {
	jsonBytes, err := json.Marshal(config)
	if err != nil {
		return err
	}

	err = fp.createBackup()
	if err != nil {
		return err
	}

	path, err := fp.configPath()
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(path, jsonBytes, 0700)
	if err != nil {
		err = fp.restoreBackup()
		if err != nil {
			return err
		}
	}

	return nil
}
