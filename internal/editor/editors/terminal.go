package options

import (
	"os"
	"os/exec"
)

// Terminal represents a plain shell session used as an editor.
type Terminal struct{}

// NewTerminal returns a new Terminal editor instance.
func NewTerminal() *Terminal {
	return &Terminal{}
}

// Name returns the name of the terminal-based editor.
func (t *Terminal) Name() string {
	return "Terminal"
}

// Open launches a shell in the project directory using $SHELL or /bin/bash.
func (t *Terminal) Open(projectPath string) error {
	shell := os.Getenv("SHELL")
	if shell == "" {
		shell = "/bin/bash"
	}

	cmd := exec.Command(shell)
	cmd.Dir = projectPath
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// IsAvailable checks whether the configured shell is available.
func (t *Terminal) IsAvailable() bool {
	shell := os.Getenv("SHELL")
	if shell == "" {
		shell = "/bin/bash"
	}

	_, err := exec.LookPath(shell)
	return err == nil
}
