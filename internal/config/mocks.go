package config

import (
	"github.com/adavidalbertson/gpair/internal/store"
)

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
			&store.InMemoryStore{},
		},
	}
}

// GetConfig returns the Config held by a MockConfigurator
// For testing purposes only
func (mc *MockConfigurator) GetConfig() Config {
	conf, _ := mc.load()
	return conf
}