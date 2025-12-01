package options

import "os/exec"

type Zed struct{}

func NewZed() *Zed {
	return &Zed{}
}

// Name returns the display name of the VSCode editor.
func (z *Zed) Name() string {
	return "Zed"
}

// Open launches Zed with the given project path.
func (v *Zed) Open(projectPath string) error {
	cmd := exec.Command("zeditor", projectPath)
	return cmd.Run()
}

// IsAvailable checks if the "zeditor" CLI command is available on the system.
func (v *Zed) IsAvailable() bool {
	_, err := exec.LookPath("zeditor")
	return err == nil
}
