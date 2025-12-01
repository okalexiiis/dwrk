package options

import (
	"os"
	"os/exec"
)

// Tmux represents a tmux session that opens inside the project directory.
type Tmux struct{}

// NewTmux returns a new Tmux instance.
func NewTmux() *Tmux {
	return &Tmux{}
}

// Name returns the display name of tmux.
func (t *Tmux) Name() string {
	return "Tmux"
}

// Open creates or attaches to a tmux session named after the project folder.
// If a session already exists, it attaches to it. Otherwise, it creates one.
func (t *Tmux) Open(projectPath string) error {
	// Extract folder name from the project path
	projectName := projectPath
	for i := len(projectPath) - 1; i >= 0; i-- {
		if projectPath[i] == '/' {
			projectName = projectPath[i+1:]
			break
		}
	}

	// Check if the session already exists
	checkCmd := exec.Command("tmux", "has-session", "-t", projectName)
	if err := checkCmd.Run(); err == nil {
		// Session exists, attach to it
		cmd := exec.Command("tmux", "attach-session", "-t", projectName)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		return cmd.Run()
	}

	// Create a new session
	cmd := exec.Command("tmux", "new-session", "-s", projectName, "-c", projectPath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// IsAvailable checks whether the tmux executable is installed.
func (t *Tmux) IsAvailable() bool {
	_, err := exec.LookPath("tmux")
	return err == nil
}
