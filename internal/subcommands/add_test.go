package subcommands

import (
	"flag"
	"reflect"
	"strings"
	"testing"

	"github.com/adavidalbertson/gpair/internal/config"
)

type mockConfigurator struct {
	pairs []config.Pair
	pair  config.Pair
	alias string
}

func (mc *mockConfigurator) GetPairs(aliases ...string) ([]config.Pair, error) {
	return mc.pairs, nil
}

func (mc *mockConfigurator) AddPair(alias string, pair config.Pair) error {
	mc.pairs = []config.Pair{pair}
	mc.pair = pair
	mc.alias = alias
	return nil
}

func TestAdd(t *testing.T) {
	tests := []struct {
		name          string
		args          string
		expectedAlias string
		expectedPair  config.Pair
	}{
		{"explicit flags", "-alias ef -email exp@flags.com -name expflags", "ef", config.Pair{Name: "expflags", Email: "exp@flags.com"}},
		{"positional args", "pa posargs pos@args.com", "pa", config.Pair{Name: "posargs", Email: "pos@args.com"}},
		{"explicit no alias", "-email exp@noalias.com -name expnoalias", "expnoalias", config.Pair{Name: "expnoalias", Email: "exp@noalias.com"}},
		{"positional no alias", "posnoalias pos@noalias.com", "posnoalias", config.Pair{Name: "posnoalias", Email: "pos@noalias.com"}},
		{"mixed", "-name mixed mix mix@mix.com", "mix", config.Pair{Name: "mixed", Email: "mix@mix.com"}},
		{"mixed no alias", "-email mix@noalias.com mixnoalias", "mixnoalias", config.Pair{Name: "mixnoalias", Email: "mix@noalias.com"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			addCmd := flag.NewFlagSet("add", flag.ExitOnError)
			addCmd.String("alias", "", "A short name for the pair")
			addCmd.String("name", "", "The git username for the pair")
			addCmd.String("email", "", "The email for the pair")

			configurator := mockConfigurator{}

			Add(strings.Split(tt.args, " "), *addCmd, &configurator)

			if configurator.alias != tt.expectedAlias {
				t.Errorf("got alias %s, want %s", configurator.alias, tt.expectedAlias)
			}

			if !reflect.DeepEqual(configurator.pair, tt.expectedPair) {
				t.Errorf("got pair %v, want %v", configurator.pair, tt.expectedPair)
			}
		})
	}
}
