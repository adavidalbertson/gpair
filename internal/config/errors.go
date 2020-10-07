package config

import (
	"fmt"
	"strings"
)

func ErrMissingCollaborator(missing []string) error {
	if len(missing) == 0 {
		return nil
	} else if len(missing) == 1 {
		return fmt.Errorf("No collaborator exists for the alias '%s'", missing[0])
	} else {
		return fmt.Errorf("No collaborators exist for aliases '%s'", strings.Join(missing, "', '"))
	}
}
