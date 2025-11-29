package open

import (
	"fmt"
	"os"

	"github.com/okalexiiis/dwrk/internal/editor"
	"github.com/okalexiiis/dwrk/internal/project"
	"github.com/okalexiiis/dwrk/pkg/utils"
	"github.com/spf13/cobra"
)

var PROJECTS_DIR = "~/Projects/"

// flags
var (
	editorFlag string
	tmuxFlag   bool
)

var OpenCmd = &cobra.Command{
	Use:   "open <nombre>",
	Short: "Abre un proyecto en el editor configurado",
	Long:  "Abre un proyecto existente en el editor de tu elección (VSCode, Neovim, etc.) o en una sesión tmux.",
	Args:  cobra.ExactArgs(1),
	Run:   runOpen,
}

func init() {
	OpenCmd.Flags().StringVarP(&editorFlag, "editor", "e", "", "Editor a usar (code, nvim, vim, nano)")
	OpenCmd.Flags().BoolVarP(&tmuxFlag, "tmux", "t", false, "Abrir en sesión tmux")
}

func runOpen(cmd *cobra.Command, args []string) {
	projectName := args[0]
	projectsDir := utils.ExpandPath(PROJECTS_DIR)

	// Verificar que el proyecto existe
	manager := project.NewManager(projectsDir)
	proj, err := manager.Get(projectName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Error: %v\n", err)
		fmt.Println("\n💡 Lista de proyectos disponibles:")
		fmt.Println("   proj list")
		os.Exit(1)
	}

	// Determinar qué editor usar
	var selectedEditor editor.Editor

	if tmuxFlag {
		selectedEditor = editor.NewTmux()
	} else if editorFlag != "" {
		// Usar el editor especificado en el flag
		selectedEditor, err = editor.GetEditor(editorFlag)
		if err != nil {
			fmt.Fprintf(os.Stderr, "❌ Error: %v\n", err)
			fmt.Println("\n📝 Editores disponibles: code, nvim, vim, nano")
			os.Exit(1)
		}
	} else {
		// Usar editor por defecto (detectado automáticamente)
		selectedEditor = editor.GetDefault()
	}

	// Abrir el proyecto
	fmt.Printf("🚀 Abriendo '%s' con %s...\n", projectName, selectedEditor.Name())

	if err := selectedEditor.Open(proj.Path); err != nil {
		fmt.Fprintf(os.Stderr, "❌ Error al abrir proyecto: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✅ Proyecto abierto exitosamente\n")
}
