package git

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/okalexiiis/dwrk/pkg/utils"
)

// Cloner provides utilities for cloning Git repositories.
type Cloner struct{}

// NewCloner creates and returns a new Cloner instance.
func NewCloner() *Cloner {
	return &Cloner{}
}

// Clone clones a Git repository into the specified directory.
//
// If the target directory does not exist, it is created automatically.
// Progress output from Git is streamed directly to stdout/stderr.
// The returned value is the absolute path of the cloned repository.
func (c *Cloner) Clone(repoURL, targetDir string) (string, error) {
	if !c.IsGitInstalled() {
		return "", fmt.Errorf("git is not installed on this system")
	}

	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return "", fmt.Errorf("failed creating target directory: %w", err)
	}

	cmd := exec.Command("git", "clone", repoURL)
	cmd.Dir = targetDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return "", utils.ParseGitError(err)
	}

	repoName := utils.ExtractRepoNameFromURL(repoURL)
	clonedPath := filepath.Join(targetDir, repoName)

	return clonedPath, nil
}

// CloneWithDepth performs a shallow clone of a Git repository.
//
// The clone will contain only the latest commit history up to the given depth.
// This is useful for performance-sensitive tasks such as CI or large repositories.
func (c *Cloner) CloneWithDepth(repoURL, targetDir string, depth int) (string, error) {
	if !c.IsGitInstalled() {
		return "", fmt.Errorf("git is not installed on this system")
	}

	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return "", fmt.Errorf("failed creating target directory: %w", err)
	}

	cmd := exec.Command("git", "clone", "--depth", fmt.Sprintf("%d", depth), repoURL)
	cmd.Dir = targetDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return "", utils.ParseGitError(err)
	}

	repoName := utils.ExtractRepoNameFromURL(repoURL)
	clonedPath := filepath.Join(targetDir, repoName)

	return clonedPath, nil
}

// CloneBranch clones a specific branch from a Git repository.
//
// If the branch does not exist, the Git command will fail and the resulting
// error will be parsed into a more readable error message.
func (c *Cloner) CloneBranch(repoURL, targetDir, branch string) (string, error) {
	if !c.IsGitInstalled() {
		return "", fmt.Errorf("git is not installed on this system")
	}

	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return "", fmt.Errorf("failed creating target directory: %w", err)
	}

	cmd := exec.Command("git", "clone", "-b", branch, repoURL)
	cmd.Dir = targetDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return "", utils.ParseGitError(err)
	}

	repoName := utils.ExtractRepoNameFromURL(repoURL)
	clonedPath := filepath.Join(targetDir, repoName)

	return clonedPath, nil
}

// IsGitInstalled checks whether the git binary is available in the system PATH.
func (c *Cloner) IsGitInstalled() bool {
	_, err := exec.LookPath("git")
	return err == nil
}
