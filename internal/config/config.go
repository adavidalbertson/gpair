package config

import (
	"fmt"
	"strings"
)

type config struct {
	Pairs map[string]Pair `json:"pairs"`
}

func newConfig() config {
	return config{
		Pairs: make(map[string]Pair),
	}
}

type Configurator interface {
	GetPairs(aliases ...string) ([]Pair, error)
	AddPair(alias string, pair Pair) error
	DeletePairs(aliases ...string) ([]string, error)
}

type configurator struct {
	provider
}

func NewConfigurator() Configurator {
	return configurator{
		&fileProvider{},
	}
}

func (c configurator) GetPairs(aliases ...string) ([]Pair, error) {
	config, err := c.load()
	if err != nil {
		return nil, err
	}

	var pairs []Pair
	var missing []string

	for _, alias := range aliases {
		if pair, ok := config.Pairs[alias]; ok {
			pairs = append(pairs, pair)
		} else {
			missing = append(missing, alias)
		}
	}

	if len(missing) == 1 {
		return pairs, fmt.Errorf("No pairing partner exists for the alias '%s'", missing[0])
	}
	if len(missing) > 1 {
		return pairs, fmt.Errorf("No pairing partners exist for aliases '%s'", strings.Join(missing, "', '"))
	}

	return pairs, nil
}

func (c configurator) AddPair(alias string, pair Pair) error {
	config, err := c.load()
	if err != nil {
		return err
	}

	config.Pairs[alias] = pair

	err = c.save(config)
	if err != nil {
		return err
	}

	return nil
}

func (c configurator) DeletePairs(aliases ...string) ([]string, error) {
	config, err := c.load()
	if err != nil {
		return nil, err
	}

	var missing []string
	var deleted []string

	for _, alias := range aliases {
		if _, exists := config.Pairs[alias]; exists {
			delete(config.Pairs, alias)
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
		return deleted, fmt.Errorf("No pairing partner exists for the alias '%s'", missing[0])
	}
	if len(missing) > 1 {
		return deleted, fmt.Errorf("No pairing partners exist for aliases '%s'", strings.Join(missing, "', '"))
	}

	return deleted, nil
}
