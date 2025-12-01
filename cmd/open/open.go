package open

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/okalexiiis/dwrk/internal/config"
	"github.com/okalexiiis/dwrk/internal/editor"
	options "github.com/okalexiiis/dwrk/internal/editor/editors"
	"github.com/okalexiiis/dwrk/internal/project"
	"github.com/spf13/cobra"
)

var (
	editorFlag string
	tmuxFlag   bool
)

var OpenCmd = &cobra.Command{
	Use:   "open <name>",
	Short: "Open a project",
	Args:  cobra.ExactArgs(1),
	Run:   runOpen,
}

func init() {
	OpenCmd.Flags().StringVarP(&editorFlag, "editor", "e", "", "Editor to use")
	OpenCmd.Flags().BoolVarP(&tmuxFlag, "tmux", "t", false, "Open the project in tmux")
}

func runOpen(cmd *cobra.Command, args []string) {
	projectName := args[0]

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading configuration: %v\n", err)
		os.Exit(1)
	}

	// Ensure the project exists
	manager := project.NewManager(cfg.ProjectsDir)
	proj, err := manager.Get(projectName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		fmt.Println("\nAvailable projects:")
		fmt.Println("  dwrk list")
		os.Exit(1)
	}

	// Determine which editor to use
	selectedEditorName := editorFlag
	if selectedEditorName == "" && !tmuxFlag {
		selectedEditorName = cfg.DefaultEditor
	}

	if tmuxFlag {
		selectedEditor := options.NewTmux()
		fmt.Printf("Opening '%s' with %s...\n", projectName, selectedEditor.Name())

		if err := selectedEditor.Open(proj.Path); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Project opened successfully")
		return
	}

	if selectedEditorName != "" && selectedEditorName != "auto" {
		selectedEditor, err := editor.GetEditor(selectedEditorName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Opening '%s' with %s...\n", projectName, selectedEditor.Name())

		if err := selectedEditor.Open(proj.Path); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Project opened successfully")
		return
	}

	// Default: open a shell session
	fmt.Printf("Opening shell in '%s'...\n", projectName)

	shell := os.Getenv("SHELL")
	if shell == "" {
		shell = "/bin/bash"
	}

	command := exec.Command(shell)
	command.Dir = proj.Path
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	command.Env = append(os.Environ(), fmt.Sprintf("PWD=%s", proj.Path))

	if err := command.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
