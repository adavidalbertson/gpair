package config

import (
	"fmt"
)

// Collaborator represents a pairing partner
type Collaborator struct {
	Alias string `json:"-"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (c Collaborator) String() string {
	return fmt.Sprintf("Co-authored-by: %s <%s>", c.Name, c.Email)
}

// NewCollaborator returns a new collaborator with the given properties
func NewCollaborator(alias, name, email string) Collaborator {
	return Collaborator{Alias: alias, Name: name, Email: email}
}

// Less returns true if a should be sorted before b, false otherwise
func Less(a, b Collaborator) bool {
	if a.Alias < b.Alias {
		return true
	}

	if a.Name < b.Name {
		return true
	}

	if a.Email < b.Email {
		return true
	}

	return false
}
