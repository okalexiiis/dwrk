package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ExtractRepoNameFromURL extracts the repository name from a Git URL.
//
// Examples:
//
//	https://github.com/user/project.git -> project
//	git@github.com:user/project.git    -> project
//	ssh://git@github.com/user/project  -> project
//
// If no name can be resolved, "cloned-repo" is returned.
func ExtractRepoNameFromURL(url string) string {
	url = strings.TrimSuffix(url, ".git")

	// Handle SSH-style URLs: git@github.com:user/repo
	if strings.Contains(url, ":") && !strings.Contains(url, "://") {
		parts := strings.Split(url, ":")
		if len(parts) >= 2 {
			url = parts[1]
		}
	}

	parts := strings.Split(url, "/")
	if len(parts) > 0 {
		last := strings.TrimSpace(parts[len(parts)-1])
		if last != "" {
			return last
		}
	}

	return "cloned-repo"
}

// ParseGitError provides a cleaner message for common Git failure patterns.
func ParseGitError(err error) error {
	errStr := err.Error()

	switch {
	case strings.Contains(errStr, "Permission denied"):
		return fmt.Errorf("permission denied: check your SSH configuration or repository access")

	case strings.Contains(errStr, "Repository not found"):
		return fmt.Errorf("repository not found: verify the URL and your access rights")

	case strings.Contains(errStr, "already exists"):
		return fmt.Errorf("target directory already exists")

	default:
		return err
	}
}

// BuildRepoURL builds an HTTPS or SSH GitHub repository URL.
func BuildRepoURL(user, repo string, https bool) string {
	if https {
		return fmt.Sprintf("https://github.com/%s/%s.git", user, repo)
	}
	return fmt.Sprintf("git@github.com:%s/%s.git", user, repo)
}

// IsGitRepo checks whether the provided directory contains a .git folder.
func IsGitRepo(path string) bool {
	gitPath := filepath.Join(path, ".git")
	info, err := os.Stat(gitPath)
	if err != nil {
		return false
	}
	return info.IsDir()
}
