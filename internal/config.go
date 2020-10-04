package internal

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

var absoluteConfigPath string

type Config struct {
	Comment bool            `json:"comment"`
	Pairs   map[string]Pair `json:"pairs"`
}

func newConfig() Config {
	return Config{
		Comment: false,
		Pairs:   make(map[string]Pair),
	}
}

func configPath() (string, error) {
	if absoluteConfigPath != "" {
		return absoluteConfigPath, nil
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
		err := os.Mkdir(configDirPath, 0755)
		if err != nil {
			return "", err
		}
	}

	configPath := filepath.Join(configDirPath, "config.json")

	if _, err = os.Stat(configPath); os.IsNotExist(err) {
		fmt.Println("Creating config")
		file, err := os.Create(configPath)
		if err != nil {
			return "", err
		}
		defer file.Close()
	}

	absoluteConfigPath = configPath
	return absoluteConfigPath, nil
}

func createBackup() error {
	path, err := configPath()
	if err != nil {
		return err
	}

	err = os.Rename(path, path+".bak")
	if err != nil {
		return err
	}

	return nil
}

func restoreBackup() error {
	path, err := configPath()
	if err != nil {
		return err
	}

	err = os.Rename(path+".bak", path)
	if err != nil {
		return err
	}

	return nil
}

func load() (Config, error) {
	path, err := configPath()
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

	var config Config
	// Since the config file is expected to be small, Marshal/Unmarshal should be adequate
	err = json.Unmarshal(jsonBytes, &config)
	if err != nil {
		return newConfig(), err
	}

	return config, nil
}

func write(config Config) error {
	jsonBytes, err := json.Marshal(config)
	if err != nil {
		return err
	}

	err = createBackup()
	if err != nil {
		return err
	}

	path, err := configPath()
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(path, jsonBytes, 0755)
	if err != nil {
		err = restoreBackup()
		if err != nil {
			return err
		}
	}

	return nil
}
