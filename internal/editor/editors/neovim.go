package options

import (
	"os"
	"os/exec"
)

// Neovim represents the Neovim text editor.
type Neovim struct{}

// NewNeovim returns a new Neovim instance.
func NewNeovim() *Neovim {
	return &Neovim{}
}

// Name returns the display name of the Neovim editor.
func (n *Neovim) Name() string {
	return "Neovim"
}

// Open launches Neovim with full terminal I/O streaming.
func (n *Neovim) Open(projectPath string) error {
	cmd := exec.Command("nvim", projectPath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// IsAvailable checks if the "nvim" executable is installed.
func (n *Neovim) IsAvailable() bool {
	_, err := exec.LookPath("nvim")
	return err == nil
}
