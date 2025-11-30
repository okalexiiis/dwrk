package editor

import (
	"fmt"
)

// Editor defines the behavior that all editors must implement.
// Each editor must provide its name, a method to open a project path,
// and a way to check if it is available on the system.
type Editor interface {
	Name() string
	Open(projectPath string) error
	IsAvailable() bool
}

// GetEditor returns an editor implementation by its name.
// Supported names: code, vscode, nvim, neovim, vim, nano.
// If the editor is not recognized or not installed, an error is returned.
func GetEditor(name string) (Editor, error) {
	var ed Editor

	switch name {
	case "code", "vscode":
		ed = NewVSCode()
	case "nvim", "neovim":
		ed = NewNeovim()
	case "vim":
		ed = NewVim()
	case "nano":
		ed = NewNano()
	default:
		return nil, fmt.Errorf("unsupported editor: %s", name)
	}

	if !ed.IsAvailable() {
		return nil, fmt.Errorf("editor '%s' is not installed on this system", name)
	}

	return ed, nil
}

// GetDefault detects and returns the first available editor based on priority.
// Priority order: VSCode > Neovim > Vim > Nano.
// If none are available, the system default editor ($EDITOR) is returned.
func GetDefault() Editor {
	editors := []Editor{
		NewVSCode(),
		NewNeovim(),
		NewVim(),
		NewNano(),
	}

	for _, ed := range editors {
		if ed.IsAvailable() {
			return ed
		}
	}

	// Fallback to the system editor when no known editors are available.
	return NewSystemEditor()
}
