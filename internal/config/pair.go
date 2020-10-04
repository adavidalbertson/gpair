package config

import (
	"fmt"
)

// Pair represents a pairing partner
type Pair struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (p Pair) String() string {
	return fmt.Sprintf("Co-authored-by: %s <%s>", p.Name, p.Email)
}