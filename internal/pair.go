package internal

import (
	"fmt"
	"strings"
)

// Pair represents a pairing partner
type Pair struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (p Pair) String() string {
	return fmt.Sprintf("Co-authored-by: %s <%s>", p.Name, p.Email)
}

func GetPairs(aliases ...string) ([]Pair, error) {
	config, err := load()
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

func AddPair(alias string, pair Pair) error {
	config, err := load()
	if err != nil {
		return err
	}

	config.Pairs[alias] = pair

	err = write(config)
	if err != nil {
		return err
	}

	return nil
}
