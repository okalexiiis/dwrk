package editor

import (
	"fmt"
)

// Editor interface que todos los editores deben implementar
type Editor interface {
	Name() string
	Open(projectPath string) error
	IsAvailable() bool
}

// GetEditor obtiene un editor específico por nombre
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
		return nil, fmt.Errorf("editor no soportado: %s", name)
	}

	if !ed.IsAvailable() {
		return nil, fmt.Errorf("el editor '%s' no está instalado en el sistema", name)
	}

	return ed, nil
}

// GetDefault detecta y retorna el editor por defecto disponible
func GetDefault() Editor {
	// Orden de prioridad: VSCode > Neovim > Vim > Nano
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

	// Fallback: abrir con $EDITOR del sistema
	return NewSystemEditor()
}
