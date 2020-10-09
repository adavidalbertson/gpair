package config

// inMemoryStore is an implementation of store that simply holds the config in memory
// For testing purposes only
type inMemoryStore struct {
	config Config
}

func (ims *inMemoryStore) load() (Config, error) {
	return ims.config, nil
}

func (ims *inMemoryStore) save(conf Config) error {
	ims.config = conf
	return nil
}

// MockConfigurator is a configurator that allows direct access to its Config
// For testing purposes only
type MockConfigurator struct {
	configurator
}

// NewMockConfigurator returns a Configurator that holds state in memory instead of writing to disk
// For testing purposes only
func NewMockConfigurator(conf Config) MockConfigurator {
	return MockConfigurator{
		configurator{
			&inMemoryStore{config: conf},
		},
	}
}

// GetConfig returns the Config held by a MockConfigurator
// For testing purposes only
func (mc *MockConfigurator) GetConfig() Config {
	conf, _ := mc.load()
	return conf
}
