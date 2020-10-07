package config

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

	return collaborators, ErrMissingCollaborator(missing)
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

	return deleted, ErrMissingCollaborator(missing)
}
