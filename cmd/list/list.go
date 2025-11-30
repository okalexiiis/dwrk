package list

import (
	"fmt"
	"os"

	"github.com/okalexiiis/dwrk/internal/config"
	"github.com/okalexiiis/dwrk/internal/project"
	"github.com/spf13/cobra"
)

var (
	showHidden bool
	filterName string
)

// ListCmd defines the `dwrk list` command.
//
// This command displays all local projects found in the configured
// projects directory. It supports filtering by name and optionally
// displaying hidden folders.
var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all local projects",
	Long:  `List all projects located in the configured projects directory.`,
	Run:   runList,
}

func init() {
	ListCmd.Flags().BoolVarP(&showHidden, "all", "a", false, "Show hidden folders")
	ListCmd.Flags().StringVarP(&filterName, "filter", "f", "", "Filter projects by name")
}

// runList executes the logic of the `list` command.
//
// Steps:
//   - Load user configuration.
//   - Initialize a project manager using the configured projects directory.
//   - Retrieve the list of projects according to CLI flags.
//   - Display the projects along with metadata (e.g., Git indicator).
//
// The command exits the program if configuration loading or project listing fails.
func runList(cmd *cobra.Command, args []string) {
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading configuration: %v\n", err)
		os.Exit(1)
	}

	manager := project.NewManager(cfg.ProjectsDir)

	projects, err := manager.List(project.ListOptions{
		ShowHidden: showHidden,
		Filter:     filterName,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error listing projects: %v\n", err)
		os.Exit(1)
	}

	if len(projects) == 0 {
		fmt.Printf("No projects found in: %s\n", cfg.ProjectsDir)
		fmt.Println("\nTip: Create a new project using:")
		fmt.Println("   dwrk new my-project")
		return
	}

	fmt.Printf("Projects in %s:\n\n", cfg.ProjectsDir)
	for i, proj := range projects {
		gitIndicator := ""
		if proj.IsGit {
			gitIndicator = " 🔗"
		}
		fmt.Printf("  %d. %s%s\n", i+1, proj.Name, gitIndicator)
	}

	fmt.Printf("\nTotal: %d project(s)\n", len(projects))
}
