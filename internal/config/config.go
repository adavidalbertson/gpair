package config

import (
	"fmt"
	"strings"
)

type config struct {
	Collaborators map[string]Collaborator `json:"collaborators"`
}

func newConfig() config {
	return config{
		Collaborators: make(map[string]Collaborator),
	}
}

type Configurator interface {
	GetCollaborators(aliases ...string) ([]Collaborator, error)
	AddCollaborator(alias string, collaborator Collaborator) error
	DeleteCollaborators(aliases ...string) ([]string, error)
}

type configurator struct {
	provider
}

func NewConfigurator() Configurator {
	return configurator{
		&fileProvider{},
	}
}

func (c configurator) GetCollaborators(aliases ...string) ([]Collaborator, error) {
	config, err := c.load()
	if err != nil {
		return nil, err
	}

	var collaborators []Collaborator
	var missing []string

	for _, alias := range aliases {
		if collaborator, ok := config.Collaborators[alias]; ok {
			collaborators = append(collaborators, collaborator)
		} else {
			missing = append(missing, alias)
		}
	}

	if len(missing) == 1 {
		return collaborators, fmt.Errorf("No collaborator exists for the alias '%s'", missing[0])
	}
	if len(missing) > 1 {
		return collaborators, fmt.Errorf("No collaborators exist for aliases '%s'", strings.Join(missing, "', '"))
	}

	return collaborators, nil
}

func (c configurator) AddCollaborator(alias string, collaborator Collaborator) error {
	config, err := c.load()
	if err != nil {
		return err
	}

	config.Collaborators[alias] = collaborator

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

	if len(missing) == 1 {
		return deleted, fmt.Errorf("No collaborator exists for the alias '%s'", missing[0])
	}
	if len(missing) > 1 {
		return deleted, fmt.Errorf("No collaborators exist for aliases '%s'", strings.Join(missing, "', '"))
	}

	return deleted, nil
}
