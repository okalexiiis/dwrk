package editor

import (
	"fmt"
	"os"
	"os/exec"
)

//
// ============================================================
// VSCode
// ============================================================
//

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

//
// ============================================================
// Neovim
// ============================================================
//

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

//
// ============================================================
// Vim
// ============================================================
//

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

//
// ============================================================
// Nano
// ============================================================
//

// Nano represents the GNU Nano editor.
type Nano struct{}

// NewNano returns a new Nano instance.
func NewNano() *Nano {
	return &Nano{}
}

// Name returns the display name of the Nano editor.
func (n *Nano) Name() string {
	return "Nano"
}

// Open launches Nano with full terminal I/O streaming.
func (n *Nano) Open(projectPath string) error {
	cmd := exec.Command("nano", projectPath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// IsAvailable checks if the "nano" executable is installed.
func (n *Nano) IsAvailable() bool {
	_, err := exec.LookPath("nano")
	return err == nil
}

//
// ============================================================
// Terminal (opens a shell inside the project directory)
// ============================================================
//

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

//
// ============================================================
// Tmux
// ============================================================
//

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

//
// ============================================================
// System Editor (fallback using $EDITOR)
// ============================================================
//

// SystemEditor represents the user-configured editor via the $EDITOR environment variable.
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
