package config

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/adavidalbertson/gpair/internal/store"
	"github.com/pkg/errors"
)

func mockStore() store.Store {
	startingConfig := populateConfig()

	configJSON, _ := json.Marshal(startingConfig)
	configBytes := []byte(configJSON)

	mockStore := &store.InMemoryStore{}
	err := mockStore.Write(configBytes)
	if err != nil {
		panic(err)
	}

	return mockStore
}

type readErrorStore struct {
	store.Store
}
func (res readErrorStore) Read() ([]byte, error) {
	return nil, errors.New("fake read error")
}
func readErrorMockStore() store.Store {
	return readErrorStore{mockStore()}
}

type writeErrorStore struct {
	store.Store
}
func (wes writeErrorStore) Write([]byte) error {
	return errors.New("fake write error")
}
func writeErrorMockStore() store.Store {
	return writeErrorStore{mockStore()}
}

type corruptedReadStore struct {
	store.Store
}
func (crs corruptedReadStore) Read() ([]byte, error) {
	return []byte("This is not a config.json!"), nil
}
func corruptedReadMockStore() store.Store {
	return corruptedReadStore{mockStore()}
}
func corruptedReadWriteErrorMockStore() store.Store {
	return corruptedReadStore{writeErrorMockStore()}
}

func populateConfig() Config {
	startingConfig := NewConfig()
	startingConfig.Collaborators["a1"] = NewCollaborator("a1", "name1", "email1")
	startingConfig.Collaborators["a2"] = NewCollaborator("a2", "name2", "email2")
	startingConfig.Collaborators["a3"] = NewCollaborator("a3", "name3", "email3")

	return startingConfig
}

func Test_configurator_GetCollaborators(t *testing.T) {
	tests := []struct {
		name    string
		aliases []string
		store   store.Store
		want    []Collaborator
		wantErr bool
	}{
		{"one alias", []string{"a2"}, mockStore(), []Collaborator{NewCollaborator("a2", "name2", "email2")}, false},
		{"two alias", []string{"a3", "a1"}, mockStore(), []Collaborator{NewCollaborator("a3", "name3", "email3"), NewCollaborator("a1", "name1", "email1")}, false},
		{"nonexistent alias", []string{"a4"}, mockStore(), nil, true},
		{"nonexistent aliases", []string{"a4", "a5"}, mockStore(), nil, true},
		{"empty args", []string{}, mockStore(), []Collaborator{NewCollaborator("a1", "name1", "email1"), NewCollaborator("a2", "name2", "email2"), NewCollaborator("a3", "name3", "email3")}, false},
		{"read error", []string{"a1"}, readErrorMockStore(), nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := configurator{
				store: tt.store,
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

	type testCase struct {
		name    string
		collaborator Collaborator
		store   store.Store
		want    Config
		wantErr bool
	}

	newTestCase := func(name string, collab Collaborator, store store.Store, wantErr bool) testCase {
		wantConfig := populateConfig()
		wantConfig.Collaborators[collab.Alias] = collab

		return testCase{
			name:    name,
			collaborator: collab,
			store:   store,
			want:    wantConfig,
			wantErr: wantErr,
		}
	}

	tests := []testCase{
		newTestCase("add new", NewCollaborator("a4", "name4", "email4"), mockStore(), false),
		newTestCase("add existing", NewCollaborator("a2", "name4", "email4"), mockStore(), false),

		// If loading fails, AddCollaborator should return an error, and load should return an empty config
		{"add read error", NewCollaborator("a4", "name4", "email4"), readErrorMockStore(), NewConfig(), true},

		// If saving fails, AddCollaborator should return an error, and the config should be unaltered
		{"add write error", NewCollaborator("a4", "name4", "email4"), writeErrorMockStore(), populateConfig(), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := configurator{
				store: tt.store,
			}
			if err := c.AddCollaborator(tt.collaborator); (err != nil) != tt.wantErr {
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
		store       store.Store
		wantDeleted []string
		wantConfig  Config
		wantErr     bool
	}

	newTestCase := func(name string, aliases []string, store store.Store, wantDeleted []string, wantErr bool) testCase {
		wantConfig := populateConfig()

		for _, alias := range aliases {
			delete(wantConfig.Collaborators, alias)
		}

		return testCase{
			name:        name,
			aliases:     aliases,
			store:       store,
			wantDeleted: wantDeleted,
			wantConfig:  wantConfig,
			wantErr:     wantErr,
		}
	}

	tests := []testCase{
		newTestCase("delete one", []string{"a2"}, mockStore(), []string{"a2"}, false),
		newTestCase("delete two", []string{"a3", "a1"}, mockStore(), []string{"a3", "a1"}, false),
		newTestCase("delete nonexistent", []string{"a4"}, mockStore(), nil, true),
		newTestCase("delete nonexistents", []string{"a4", "a5"}, mockStore(), nil, true),
		newTestCase("delete mixed", []string{"a2", "a4"}, mockStore(), []string{"a2"}, true),

		// if loading fails, DeleteCollaborators should return nil and an error, and load will return an empty config
		{"delete read error", []string{"a2"}, readErrorMockStore(), nil, NewConfig(), true},

		// if saving fails, DeleteCollaborators should return nil and an error, and the config should be unaltered
		{"delete write error", []string{"a2"}, writeErrorMockStore(), nil, populateConfig(), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := configurator{
				store: tt.store,
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

func Test_configurator_load(t *testing.T) {
	tests := []struct {
		name    string
		store   store.Store
		want    Config
		wantErr bool
	}{
		{"happy path", mockStore(), populateConfig(), false},
		{"empty file", &store.InMemoryStore{}, NewConfig(), false},
		{"read error", readErrorMockStore(), NewConfig(), true},

		// If load fails json parsing, it returns an empty config and no error, and writes an empty config
		{"corrupted file", corruptedReadMockStore(), NewConfig(), false},

		// If the write of the empty config fails, *then* it returns an error
		{"corrupted file and write error", corruptedReadWriteErrorMockStore(), NewConfig(), true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := configurator{
				store: tt.store,
			}
			got, err := c.load()
			if (err != nil) != tt.wantErr {
				t.Errorf("configurator.load() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("configurator.load() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_configurator_save(t *testing.T) {
	type args struct {
		config Config
	}
	tests := []struct {
		name       string
		store      store.Store
		args       args
		wantConfig Config
		wantErr    bool
	}{
		{"happy path", mockStore(), args{populateConfig()}, populateConfig(), false},
		{"empty config", mockStore(), args{NewConfig()}, NewConfig(), false},
		{"write error", writeErrorMockStore(), args{populateConfig()}, populateConfig(), true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := configurator{
				store: tt.store,
			}
			if err := c.save(tt.args.config); (err != nil) != tt.wantErr {
				t.Errorf("configurator.save() error = %v, wantErr %v", err, tt.wantErr)
			}

			gotConfig, _ := c.load()
			if !reflect.DeepEqual(gotConfig, tt.wantConfig) {
				t.Errorf("configurator.load() = %v, want %v", gotConfig, tt.wantConfig)
			}
		})
	}
}
