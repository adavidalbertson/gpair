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

type Configurator struct {
	provider
}

func NewConfigurator() Configurator {
	return Configurator{
		&fileProvider{},
	}
}

func (c Configurator) GetPairs(aliases ...string) ([]Pair, error) {
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

func (c Configurator) AddPair(alias string, pair Pair) error {
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
