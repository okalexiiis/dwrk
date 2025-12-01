package options

import (
	"fmt"
	"os"
	"os/exec"
)

type SystemEditor struct{}

// NewSystemEditor returns a new SystemEditor instance.
func NewSystemEditor() *SystemEditor {
	return &SystemEditor{}
}

// Name returns the name of the system editor from $EDITOR.
// If not defined, a generic fallback string is returned.
func (s *SystemEditor) Name() string {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		return "System Editor"
	}
	return editor
}

// Open launches the editor defined in $EDITOR, if available.
func (s *SystemEditor) Open(projectPath string) error {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		return fmt.Errorf("no system editor configured. Set $EDITOR or specify --editor")
	}

	cmd := exec.Command(editor, projectPath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// IsAvailable checks if the editor in $EDITOR exists on the system.
func (s *SystemEditor) IsAvailable() bool {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		return false
	}
	_, err := exec.LookPath(editor)
	return err == nil
}
