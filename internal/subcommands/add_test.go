package subcommands

import (
	"flag"
	"reflect"
	"strings"
	"testing"

	"github.com/adavidalbertson/gpair/internal/config"
)

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
			AddCmd.String("alias", "", "A short name for the collaborator")
			AddCmd.String("name", "", "The git username for the collaborator")
			AddCmd.String("email", "", "The email for the collaborator")

			gotAlias, gotName, gotEmail, err := parseAddArgs(strings.Split(tt.args, " "))

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

			configurator := config.NewMockConfigurator(config.NewConfig())

			err := add(tt.args.alias, tt.args.name, tt.args.email, configurator)

			if err != nil != tt.wantErr {
				t.Errorf("Add() error = %v, wantErr %v", err, tt.wantErr)
			}

			if err == nil {
				wantConfig := config.NewConfig()
				wantConfig.Collaborators[tt.args.alias] = config.NewCollaborator(tt.args.alias, tt.args.name, tt.args.email)

				gotConfig := configurator.GetConfig()

				if !reflect.DeepEqual(gotConfig, wantConfig) {
					t.Errorf("config = %v, want %v", gotConfig, wantConfig)
				}
			}
		})
	}
}
