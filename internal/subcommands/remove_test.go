package subcommands

import (
	"strings"
	"reflect"
	"testing"
)

func TestParseRemoveArgs(t *testing.T) {
	tests := []struct {
		name        string
		args        string
		wantAliases []string
		wantErr     bool
	}{
		{"happy path", "da sn", []string{"da", "sn"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotAliases, err := parseRemoveArgs(strings.Split(tt.args, " "))
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseRemoveArgs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotAliases, tt.wantAliases) {
				t.Errorf("ParseRemoveArgs() = %v, want %v", gotAliases, tt.wantAliases)
			}
		})
	}
}
