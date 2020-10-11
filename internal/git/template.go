package git

import (
	"github.com/adavidalbertson/gpair/internal/config"
	"github.com/adavidalbertson/gpair/internal/store"
)

// CreateTemplate saves a git commit template containing the paired collaborators
func CreateTemplate(repoName string, coauthors ...config.Collaborator) (string, error) {
	store, err := store.NewFileStore(repoName + "-template.txt", store.HOME, ".gpair")
	if err != nil {
		return "", err
	}

	template := "\n\n# Co-author trailer provided by gpair\n\n"

	for _, coauthor := range coauthors {
		template += coauthor.String() + "\n"
	}

	err = store.Write([]byte(template))
	if err != nil {
		return "", err
	}

	return store.GetPath(), nil
}