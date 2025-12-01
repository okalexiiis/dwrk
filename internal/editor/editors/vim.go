package options

import (
	"os"
	"os/exec"
)

// Vim represents the Vim text editor.
type Vim struct{}

// NewVim returns a new Vim instance.
func NewVim() *Vim {
	return &Vim{}
}

// Name returns the display name of the Vim editor.
func (v *Vim) Name() string {
	return "Vim"
}

// Open launches Vim with full terminal I/O streaming.
func (v *Vim) Open(projectPath string) error {
	cmd := exec.Command("vim", projectPath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// IsAvailable checks if the "vim" binary is installed.
func (v *Vim) IsAvailable() bool {
	_, err := exec.LookPath("vim")
	return err == nil
}
