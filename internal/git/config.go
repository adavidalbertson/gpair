package git

import (
	"strings"
	"path/filepath"
	"os/exec"
)

// IsInstalled returns true if git is on the user's path
func IsInstalled() bool {
	_, err := exec.LookPath("git")
	return err == nil
}

// GetRepoName returns the name of the git repo where gpair was executed
func GetRepoName() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	repoPathBytes, err := cmd.Output()
	if err != nil {
		return "", err
	}

	_, repoName := filepath.Split(string(repoPathBytes))

	return strings.TrimSpace(repoName), nil
}

// IsCustomTemplate returns true if git is already configured with a template not made by gpair
func IsCustomTemplate() (bool, error) {
	cmd := exec.Command("git", "config", "--get", "--null", "commit.template")
	templatePathBytes, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			if exitErr.ExitCode() == 1 {
				// git config exits with code 1 if the config is not set
				return false, nil
			}
		}

		return false, err
	}

	templatePath := string(templatePathBytes)

	return len(templatePath) > 0 && !strings.Contains(templatePath, ".gpair"), nil
}

// SetTemplate sets the current repo's git config commit.template to the provided filepath
func SetTemplate(templatePath string, global bool) error {
	cmd := exec.Command("git", gitConfig(global, "commit.template", templatePath)...)
	return cmd.Run()
}

// UnsetTemplate unsets the current repo's git config commit.template
func UnsetTemplate(global bool) error {
	cmd := exec.Command("git", gitConfig(global, "--unset", "commit.template")...)
	err := cmd.Run()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			if exitErr.ExitCode() == 5 {
				// git config exits with code 5 if the config does not exist
				// or if multiple config lines match
				return nil
			}
		}
	}
	
	return err
}

func gitConfig(global bool, args ...string) []string {
	cmdString := []string{"config"}
	if global {
		cmdString = append(cmdString, "--global")
	}
	return append(cmdString, args...)
}