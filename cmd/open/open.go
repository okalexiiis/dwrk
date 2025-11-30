package open

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/okalexiiis/dwrk/internal/config"
	"github.com/okalexiiis/dwrk/internal/editor"
	"github.com/okalexiiis/dwrk/internal/project"
	"github.com/spf13/cobra"
)

var (
	editorFlag string
	tmuxFlag   bool
)

var OpenCmd = &cobra.Command{
	Use:   "open <nombre>",
	Short: "Abre un proyecto",
	Args:  cobra.ExactArgs(1),
	Run:   runOpen,
}

func init() {
	OpenCmd.Flags().StringVarP(&editorFlag, "editor", "e", "", "Editor a usar")
	OpenCmd.Flags().BoolVarP(&tmuxFlag, "tmux", "t", false, "Abrir en tmux")
}

func runOpen(cmd *cobra.Command, args []string) {
	projectName := args[0]

	// Cargar configuración
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Error cargando configuración: %v\n", err)
		os.Exit(1)
	}

	// Verificar que el proyecto existe
	manager := project.NewManager(cfg.ProjectsDir)
	proj, err := manager.Get(projectName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Error: %v\n", err)
		fmt.Println("\n💡 Lista de proyectos disponibles:")
		fmt.Println("   dwrk list")
		os.Exit(1)
	}

	// Determinar editor a usar
	selectedEditorName := editorFlag
	if selectedEditorName == "" && !tmuxFlag {
		selectedEditorName = cfg.DefaultEditor
	}

	if tmuxFlag {
		selectedEditor := editor.NewTmux()
		fmt.Printf("🚀 Abriendo '%s' con %s...\n", projectName, selectedEditor.Name())
		if err := selectedEditor.Open(proj.Path); err != nil {
			fmt.Fprintf(os.Stderr, "❌ Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("✅ Proyecto abierto exitosamente\n")
	} else if selectedEditorName != "" && selectedEditorName != "auto" {
		selectedEditor, err := editor.GetEditor(selectedEditorName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "❌ Error: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("🚀 Abriendo '%s' con %s...\n", projectName, selectedEditor.Name())
		if err := selectedEditor.Open(proj.Path); err != nil {
			fmt.Fprintf(os.Stderr, "❌ Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("✅ Proyecto abierto exitosamente\n")
	} else {
		// Default: abrir shell
		fmt.Printf("📁 Abriendo shell en '%s'...\n", projectName)
		shell := os.Getenv("SHELL")
		if shell == "" {
			shell = "/bin/bash"
		}
		cmd := exec.Command(shell)
		cmd.Dir = proj.Path
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Env = append(os.Environ(), fmt.Sprintf("PWD=%s", proj.Path))
		if err := cmd.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "❌ Error: %v\n", err)
			os.Exit(1)
		}
	}
}
