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
		pair, ok := config.Pairs[alias]
		if ok {
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
