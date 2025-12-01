package new

import (
	"fmt"
	"os"

	"github.com/okalexiiis/dwrk/internal/config"
	"github.com/okalexiiis/dwrk/internal/project"
	"github.com/okalexiiis/dwrk/pkg/utils"
	"github.com/spf13/cobra"
)

// flags
var (
	git      bool
	template string
)

var NewCmd = &cobra.Command{
	Use:   "new <name>",
	Short: "Create a new local project",
	Long:  "Creates a new local project in the configured directory and optionally initializes a Git repository.",
	Args:  cobra.ExactArgs(1),
	Run:   runNew,
}

func init() {
	NewCmd.Flags().BoolVarP(&git, "git", "g", false, "Initialize a Git repository")
	NewCmd.Flags().StringVarP(&template, "template", "t", "", "Create a Project with a template")
}

func runNew(cmd *cobra.Command, args []string) {
	// Load config
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading configuration: %v\n", err)
		os.Exit(1)
	}

	projectName := args[0]
	projectsDir := utils.ExpandPath(cfg.ProjectsDir)

	// Create project manager
	manager := project.NewManager(projectsDir)

	// Attempt to create the project
	createdProject, err := manager.Create(projectName, project.CreateOptions{
		InitGit:  git,
		Template: template,
	})

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Success message
	fmt.Printf("Project created successfully: %s\n", createdProject.Name)
	fmt.Printf("Location: %s\n", createdProject.Path)

	if createdProject.IsGit {
		fmt.Println("Git repository initialized")
	}

	fmt.Println("\nTo open the project:")
	fmt.Printf("  dwrk open %s\n", projectName)
}
