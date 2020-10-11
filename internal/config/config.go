package config

// Config is the persisted config for gpair, including a dictionary of collaborators
type Config struct {
	Collaborators map[string]Collaborator `json:"collaborators"`
}

// NewConfig returns an empty Config
func NewConfig() Config {
	return Config{
		Collaborators: make(map[string]Collaborator),
	}
}
