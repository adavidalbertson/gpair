package config

import (
	"encoding/json"

	"github.com/adavidalbertson/gpair/internal"
	"github.com/adavidalbertson/gpair/internal/store"
)

// Configurator is an abstraction that allows operations to a persisted config
type Configurator interface {
	GetCollaborators(aliases ...string) ([]Collaborator, error)
	AddCollaborator(collaborator Collaborator) error
	DeleteCollaborators(aliases ...string) ([]string, error)
}

type configurator struct {
	store store.Store
}

// NewConfigurator returns a configurator that persists config to disk
func NewConfigurator() (Configurator, error) {
	store, err := store.NewFileStore("config.json", store.HOME, ".gpair")
	if err != nil {
		return nil, err
	}

	return configurator{store}, nil
}

func (c configurator) load() (Config, error) {
	jsonBytes, err := c.store.Read()
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
		err = c.store.Write([]byte{})
		if err != nil {
			return NewConfig(), err
		}
		path := c.store.GetPath()
		internal.PrintVerbose("Failed to parse config file at %s. Starting a new config file. Original has been moved to %s.bak", path, path)
		return NewConfig(), nil
	}
	
	for alias, collab := range config.Collaborators {
		collab.Alias = alias
		config.Collaborators[alias] = collab
	}

	return config, nil
}

func (c configurator) save(config Config) error {
	jsonBytes, err := json.Marshal(config)
	if err != nil {
		return err
	}

	err = c.store.Write(jsonBytes)
	if err != nil {
		return err
	}

	return nil
}

func (c configurator) GetCollaborators(aliases ...string) ([]Collaborator, error) {
	config, err := c.load()
	if err != nil {
		return nil, err
	}

	var collaborators []Collaborator
	var missing []string

	if len(aliases) == 0 {
		for _, collab := range config.Collaborators {
			collaborators = append(collaborators, collab)
		}
	}

	for _, alias := range aliases {
		if collaborator, ok := config.Collaborators[alias]; ok {
			collaborators = append(collaborators, collaborator)
		} else {
			missing = append(missing, alias)
		}
	}

	return collaborators, ErrMissingCollaborator(missing)
}

func (c configurator) AddCollaborator(collab Collaborator) error {
	config, err := c.load()
	if err != nil {
		return err
	}

	if len(collab.Alias) == 0 {
		collab.Alias = collab.Name
	}

	config.Collaborators[collab.Alias] = collab

	err = c.save(config)
	if err != nil {
		return err
	}

	return nil
}

func (c configurator) DeleteCollaborators(aliases ...string) ([]string, error) {
	config, err := c.load()
	if err != nil {
		return nil, err
	}

	var missing []string
	var deleted []string

	for _, alias := range aliases {
		if _, exists := config.Collaborators[alias]; exists {
			delete(config.Collaborators, alias)
			deleted = append(deleted, alias)
		} else {
			missing = append(missing, alias)
		}
	}

	err = c.save(config)
	if err != nil {
		return nil, err
	}

	return deleted, ErrMissingCollaborator(missing)
}
