package options

import "os/exec"

// VSCode represents the Visual Studio Code editor.
type VSCode struct{}

// NewVSCode returns a new VSCode instance.
func NewVSCode() *VSCode {
	return &VSCode{}
}

// Name returns the display name of the VSCode editor.
func (v *VSCode) Name() string {
	return "VSCode"
}

// Open launches VSCode with the given project path.
func (v *VSCode) Open(projectPath string) error {
	cmd := exec.Command("code", projectPath)
	return cmd.Run()
}

// IsAvailable checks if the "code" CLI command is available on the system.
func (v *VSCode) IsAvailable() bool {
	_, err := exec.LookPath("code")
	return err == nil
}
