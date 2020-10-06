package subcommands

import (
	"flag"
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

func TestParseAddArgs(t *testing.T) {
	tests := []struct {
		name      string
		args      string
		wantAlias string
		wantName  string
		wantEmail string
	}{
		{"explicit flags", "-alias ef -email exp@flags.com -name expflags", "ef", "expflags", "exp@flags.com"},
		{"positional args", "pa posargs pos@args.com", "pa", "posargs", "pos@args.com"},
		{"explicit no alias", "-email exp@noalias.com -name expnoalias", "expnoalias", "expnoalias", "exp@noalias.com"},
		{"positional no alias", "posnoalias pos@noalias.com", "posnoalias", "posnoalias", "pos@noalias.com"},
		{"mixed", "-name mixed mix mix@mix.com", "mix", "mixed", "mix@mix.com"},
		{"mixed no alias", "-email mix@noalias.com mixnoalias", "mixnoalias", "mixnoalias", "mix@noalias.com"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			AddCmd = *flag.NewFlagSet("add", flag.ExitOnError)
			AddCmd.String("alias", "", "A short name for the pair")
			AddCmd.String("name", "", "The git username for the pair")
			AddCmd.String("email", "", "The email for the pair")

			gotAlias, gotName, gotEmail, err := ParseAddArgs(strings.Split(tt.args, " "))

			if err != nil {
				t.Errorf("subcommands.ParseAddArgs() error = %v, wantErr false", err)
			}

			if gotAlias != tt.wantAlias {
				t.Errorf("got alias %s, want %s", gotAlias, tt.wantAlias)
			}

			if gotName != tt.wantName {
				t.Errorf("got name %s, want %s", gotName, tt.wantName)
			}

			if gotEmail != tt.wantEmail {
				t.Errorf("got name %s, want %s", gotEmail, tt.wantEmail)
			}
		})
	}
}

func TestAdd(t *testing.T) {
	type args struct {
		alias string
		name  string
		email string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"happy path", args{"happy", "happypath", "happyemail"}, false},
		{"missing name", args{"sad", "", "sademail"}, true},
		{"missing email", args{"sad", "sadpath", ""}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			configurator := mockConfigurator{}

			err := Add(tt.args.alias, tt.args.name, tt.args.email, &configurator)

			if err != nil != tt.wantErr {
				t.Errorf("Add() error = %v, wantErr %v", err, tt.wantErr)
			}

			if err == nil {
				if configurator.alias != tt.args.alias {
					t.Errorf("got alias %s, want %s", configurator.alias, tt.args.alias)
				}

				if configurator.pair.Name != tt.args.name {
					t.Errorf("got name %s, want %s", configurator.pair.Name, tt.args.name)
				}

				if configurator.pair.Email != tt.args.email {
					t.Errorf("got alias %s, want %s", configurator.pair.Email, tt.args.email)
				}
			}
		})
	}
}
