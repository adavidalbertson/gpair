package config

import (
	"fmt"
)

// Collaborator represents a pairing partner
type Collaborator struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (c Collaborator) String() string {
	return fmt.Sprintf("Co-authored-by: %s <%s>", c.Name, c.Email)
}

// NewCollaborator returns a new collaborator with the given properties
func NewCollaborator(name, email string) Collaborator {
	return Collaborator{Name: name, Email: email}
}