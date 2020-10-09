package config

// Config is the persisted config for gpair, including a dictionary of collaborators
type Config struct {
	Collaborators map[string]Collaborator `json:"collaborators"`
}

// NewConfig returns an empty Config
func NewConfig() Config {
	return Config{
		Collaborators: make(map[string]Collaborator),
	}
}

// Configurator is an abstraction that allows operations to a persisted config
type Configurator interface {
	GetCollaborators(aliases ...string) ([]Collaborator, error)
	AddCollaborator(alias string, collaborator Collaborator) error
	DeleteCollaborators(aliases ...string) ([]string, error)
}

type configurator struct {
	store
}

// NewConfigurator returns a configurator that persists config to disk
func NewConfigurator() Configurator {
	return configurator{
		&fileStore{},
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
