package config

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/adavidalbertson/gpair/internal/store"
)

var mockStore *store.InMemoryStore

func setUp() {
	startingConfig := populateConfig()

	configJSON, _ := json.Marshal(startingConfig)
	configBytes := []byte(configJSON)

	mockStore = &store.InMemoryStore{}
	mockStore.Write(configBytes)
}

func populateConfig() Config {
	startingConfig := NewConfig()
	startingConfig.Collaborators["a1"] = NewCollaborator("name1", "email1")
	startingConfig.Collaborators["a2"] = NewCollaborator("name2", "email2")
	startingConfig.Collaborators["a3"] = NewCollaborator("name3", "email3")

	return startingConfig
}

func Test_configurator_GetCollaborators(t *testing.T) {
	setUp()

	tests := []struct {
		name    string
		aliases []string
		want    []Collaborator
		wantErr bool
	}{
		{"one alias", []string{"a2"}, []Collaborator{NewCollaborator("name2", "email2")}, false},
		{"two alias", []string{"a3", "a1"}, []Collaborator{NewCollaborator("name3", "email3"), NewCollaborator("name1", "email1")}, false},
		{"nonexistent alias", []string{"a4"}, nil, true},
		{"nonexistent aliases", []string{"a4", "a5"}, nil, true},
		{"empty args", []string{}, nil, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := configurator{
				store: mockStore,
			}
			got, err := c.GetCollaborators(tt.aliases...)
			if (err != nil) != tt.wantErr {
				t.Errorf("configurator.GetCollaborators() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("configurator.GetCollaborators() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_configurator_AddCollaborator(t *testing.T) {

	type args struct {
		alias        string
		collaborator Collaborator
	}
	type testCase struct {
		name    string
		args    args
		want    Config
		wantErr bool
	}

	newTestCase := func(name, alias string, collab Collaborator, wantErr bool) testCase {
		wantConfig := populateConfig()
		wantConfig.Collaborators[alias] = collab

		return testCase{
			name:    name,
			args:    args{alias, collab},
			want:    wantConfig,
			wantErr: wantErr,
		}
	}

	tests := []testCase{
		newTestCase("add new", "a4", NewCollaborator("name4", "email4"), false),
		newTestCase("add existing", "a2", NewCollaborator("name4", "email4"), false),
	}

	for _, tt := range tests {

		setUp()

		t.Run(tt.name, func(t *testing.T) {
			c := configurator{
				store: mockStore,
			}
			if err := c.AddCollaborator(tt.args.alias, tt.args.collaborator); (err != nil) != tt.wantErr {
				t.Errorf("configurator.AddCollaborator() error = %v, wantErr %v", err, tt.wantErr)
			}

			got, _ := c.load()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("configurator.load() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_configurator_DeleteCollaborators(t *testing.T) {
	type testCase struct {
		name        string
		aliases     []string
		wantDeleted []string
		wantConfig  Config
		wantErr     bool
	}

	newTestCase := func(name string, aliases, wantDeleted []string, wantErr bool) testCase {
		wantConfig := populateConfig()

		for _, alias := range aliases {
			delete(wantConfig.Collaborators, alias)
		}

		return testCase{
			name:        name,
			aliases:     aliases,
			wantDeleted: wantDeleted,
			wantConfig:  wantConfig,
			wantErr:     wantErr,
		}
	}

	tests := []testCase{
		newTestCase("delete one", []string{"a2"}, []string{"a2"}, false),
		newTestCase("delete two", []string{"a3", "a1"}, []string{"a3", "a1"}, false),
		newTestCase("delete nonexistent", []string{"a4"}, nil, true),
		newTestCase("delete nonexistents", []string{"a4", "a5"}, nil, true),
		newTestCase("delete mixed", []string{"a2", "a4"}, []string{"a2"}, true),
	}

	for _, tt := range tests {

		setUp()

		t.Run(tt.name, func(t *testing.T) {
			c := configurator{
				store: mockStore,
			}
			gotDeleted, err := c.DeleteCollaborators(tt.aliases...)
			if (err != nil) != tt.wantErr {
				t.Errorf("configurator.DeleteCollaborators() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotDeleted, tt.wantDeleted) {
				t.Errorf("configurator.DeleteCollaborators() = %v, want %v", gotDeleted, tt.wantDeleted)
			}

			gotConfig, _ := c.load()
			if !reflect.DeepEqual(gotConfig, tt.wantConfig) {
				t.Errorf("configurator.load() = %v, want %v", gotConfig, tt.wantConfig)
			}
		})
	}
}
