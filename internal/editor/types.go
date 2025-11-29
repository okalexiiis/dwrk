package editor

import (
	"fmt"
	"os"
	"os/exec"
)

// ============================================
// VSCode
// ============================================

type VSCode struct{}

func NewVSCode() *VSCode {
	return &VSCode{}
}

func (v *VSCode) Name() string {
	return "VSCode"
}

func (v *VSCode) Open(projectPath string) error {
	cmd := exec.Command("code", projectPath)
	return cmd.Run()
}

func (v *VSCode) IsAvailable() bool {
	_, err := exec.LookPath("code")
	return err == nil
}

// ============================================
// Neovim
// ============================================

type Neovim struct{}

func NewNeovim() *Neovim {
	return &Neovim{}
}

func (n *Neovim) Name() string {
	return "Neovim"
}

func (n *Neovim) Open(projectPath string) error {
	cmd := exec.Command("nvim", projectPath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (n *Neovim) IsAvailable() bool {
	_, err := exec.LookPath("nvim")
	return err == nil
}

// ============================================
// Vim
// ============================================

type Vim struct{}

func NewVim() *Vim {
	return &Vim{}
}

func (v *Vim) Name() string {
	return "Vim"
}

func (v *Vim) Open(projectPath string) error {
	cmd := exec.Command("vim", projectPath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (v *Vim) IsAvailable() bool {
	_, err := exec.LookPath("vim")
	return err == nil
}

// ============================================
// Nano
// ============================================

type Nano struct{}

func NewNano() *Nano {
	return &Nano{}
}

func (n *Nano) Name() string {
	return "Nano"
}

func (n *Nano) Open(projectPath string) error {
	cmd := exec.Command("nano", projectPath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (n *Nano) IsAvailable() bool {
	_, err := exec.LookPath("nano")
	return err == nil
}

// =======================================
// Terminal
// =======================================
type Terminal struct{}

func NewTerminal() *Terminal {
	return &Terminal{}
}

func (t *Terminal) Name() string {
	return "Terminal"
}

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

func (t *Terminal) IsAvailable() bool {
	shell := os.Getenv("SHELL")
	if shell == "" {
		shell = "/bin/bash"
	}

	_, err := exec.LookPath(shell)
	return err == nil
}

// ============================================
// Tmux
// ============================================

type Tmux struct{}

func NewTmux() *Tmux {
	return &Tmux{}
}

func (t *Tmux) Name() string {
	return "Tmux"
}

func (t *Tmux) Open(projectPath string) error {
	// Obtener el nombre del proyecto del path
	projectName := projectPath[len(projectPath)-1:]
	for i := len(projectPath) - 1; i >= 0; i-- {
		if projectPath[i] == '/' {
			projectName = projectPath[i+1:]
			break
		}
	}

	// Verificar si ya existe una sesión con este nombre
	checkCmd := exec.Command("tmux", "has-session", "-t", projectName)
	if err := checkCmd.Run(); err == nil {
		// La sesión ya existe, adjuntarse a ella
		cmd := exec.Command("tmux", "attach-session", "-t", projectName)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		return cmd.Run()
	}

	// Crear nueva sesión
	cmd := exec.Command("tmux", "new-session", "-s", projectName, "-c", projectPath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (t *Tmux) IsAvailable() bool {
	_, err := exec.LookPath("tmux")
	return err == nil
}

// ============================================
// System Editor (fallback usando $EDITOR)
// ============================================

type SystemEditor struct{}

func NewSystemEditor() *SystemEditor {
	return &SystemEditor{}
}

func (s *SystemEditor) Name() string {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		return "Editor del sistema"
	}
	return editor
}

func (s *SystemEditor) Open(projectPath string) error {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		return fmt.Errorf("no hay editor configurado. Define $EDITOR o usa --editor")
	}

	cmd := exec.Command(editor, projectPath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (s *SystemEditor) IsAvailable() bool {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		return false
	}
	_, err := exec.LookPath(editor)
	return err == nil
}
