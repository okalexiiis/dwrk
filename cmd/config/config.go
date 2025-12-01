package config

import (
	"fmt"
	"os"

	"github.com/okalexiiis/dwrk/internal/config"
	"github.com/spf13/cobra"
)

var ConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "Gestiona la configuración de dwrk",
	Long: `Gestiona la configuración de dwrk.
	// Subcommands:
//   set     Set a configuration value
//   get     Get a configuration value
//   list    List all configuration values
//   path    Show the configuration file path
//   reset   Reset all configuration to defaults`,
}

// Set command
var setCmd = &cobra.Command{
	Use:   "set <key> <value>",
	Short: "Set a configuration value",
	Long: `Set a configuration key.

Available keys:
  projects_dir      Base directory for projects
  editor            Default editor (auto, code, nvim, vim, nano, terminal)
  github_username   GitHub username
  use_ssh           Use SSH for Git operations (true/false)

Examples:
  dwrk config set projects_dir ~/Dev
  dwrk config set editor code
  dwrk config set github_username myuser
  dwrk config set use_ssh false`,
	Args: cobra.ExactArgs(2),
	Run:  runSet,
}

// Get command
var getCmd = &cobra.Command{
	Use:   "get <key>",
	Short: "Get a configuration value",
	Long: `Retrieve a configuration value.

Examples:
  dwrk config get projects_dir
  dwrk config get editor`,
	Args: cobra.ExactArgs(1),
	Run:  runGet,
}

// List command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all configuration values",
	Run:   runList,
}

// Path command
var pathCmd = &cobra.Command{
	Use:   "path",
	Short: "Show the configuration file path",
	Run:   runPath,
}

// Reset command
var resetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Reset configuration to default values",
	Run:   runReset,
}

func init() {
	ConfigCmd.AddCommand(setCmd)
	ConfigCmd.AddCommand(getCmd)
	ConfigCmd.AddCommand(listCmd)
	ConfigCmd.AddCommand(pathCmd)
	ConfigCmd.AddCommand(resetCmd)
}

func runSet(cmd *cobra.Command, args []string) {
	key := args[0]
	value := args[1]

	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading configuration: %v\n", err)
		os.Exit(1)
	}

	if err := cfg.Set(key, value); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Configuration updated: %s = %s\n", key, value)
}

func runGet(cmd *cobra.Command, args []string) {
	key := args[0]

	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading configuration: %v\n", err)
		os.Exit(1)
	}

	value, err := cfg.Get(key)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(value)
}

func runList(cmd *cobra.Command, args []string) {
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading configuration: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Current configuration:")
	fmt.Println()
	fmt.Printf("  projects_dir:     %s\n", cfg.ProjectsDir)
	fmt.Printf("  templates_dir:    %s\n", cfg.TemplatesDir)
	fmt.Printf("  default_editor:   %s\n", cfg.DefaultEditor)
	fmt.Printf("  github_username:  %s\n", cfg.GitHubUsername)
	fmt.Printf("  use_ssh:          %v\n", cfg.UseSSH)
	fmt.Println()
	fmt.Printf("Configuration file: %s\n", config.GetConfigPath())
}

func runPath(cmd *cobra.Command, args []string) {
	fmt.Println(config.GetConfigPath())
}

func runReset(cmd *cobra.Command, args []string) {
	fmt.Print("Are you sure you want to reset the configuration? [y/N]: ")

	var response string
	fmt.Scanln(&response)

	if response != "y" && response != "Y" && response != "yes" {
		fmt.Println("Operation cancelled")
		return
	}

	cfg := config.Default()
	if err := cfg.Save(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Configuration reset to default values")
}
