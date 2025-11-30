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

Subcomandos:
  set     Establece un valor de configuración
  get     Obtiene un valor de configuración
  list    Lista toda la configuración
  path    Muestra la ruta del archivo de configuración
  reset   Restablece la configuración a valores por defecto`,
}

// Set command
var setCmd = &cobra.Command{
	Use:   "set <clave> <valor>",
	Short: "Establece un valor de configuración",
	Long: `Establece un valor de configuración.

Claves disponibles:
  projects_dir      Directorio base de proyectos
  editor            Editor por defecto (auto, code, nvim, vim, nano, terminal)
  github_username   Usuario de GitHub
  use_ssh           Usar SSH para git (true/false)

Ejemplos:
  dwrk config set projects_dir ~/Dev
  dwrk config set editor code
  dwrk config set github_username miusuario
  dwrk config set use_ssh false`,
	Args: cobra.ExactArgs(2),
	Run:  runSet,
}

// Get command
var getCmd = &cobra.Command{
	Use:   "get <clave>",
	Short: "Obtiene un valor de configuración",
	Long: `Obtiene un valor de configuración.

Ejemplos:
  dwrk config get projects_dir
  dwrk config get editor`,
	Args: cobra.ExactArgs(1),
	Run:  runGet,
}

// List command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Lista toda la configuración",
	Run:   runList,
}

// Path command
var pathCmd = &cobra.Command{
	Use:   "path",
	Short: "Muestra la ruta del archivo de configuración",
	Run:   runPath,
}

// Reset command
var resetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Restablece la configuración a valores por defecto",
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
		fmt.Fprintf(os.Stderr, "❌ Error cargando configuración: %v\n", err)
		os.Exit(1)
	}

	if err := cfg.Set(key, value); err != nil {
		fmt.Fprintf(os.Stderr, "❌ Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✅ Configuración actualizada: %s = %s\n", key, value)
}

func runGet(cmd *cobra.Command, args []string) {
	key := args[0]

	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Error cargando configuración: %v\n", err)
		os.Exit(1)
	}

	value, err := cfg.Get(key)
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(value)
}

func runList(cmd *cobra.Command, args []string) {
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Error cargando configuración: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("⚙️  Configuración actual:")
	fmt.Println()
	fmt.Printf("  projects_dir:     %s\n", cfg.ProjectsDir)
	fmt.Printf("  default_editor:   %s\n", cfg.DefaultEditor)
	fmt.Printf("  github_username:  %s\n", cfg.GitHubUsername)
	fmt.Printf("  use_ssh:          %v\n", cfg.UseSSH)
	fmt.Println()
	fmt.Printf("📄 Archivo de configuración: %s\n", config.GetConfigPath())
}

func runPath(cmd *cobra.Command, args []string) {
	fmt.Println(config.GetConfigPath())
}

func runReset(cmd *cobra.Command, args []string) {
	fmt.Print("⚠️  ¿Estás seguro de restablecer la configuración? [y/N]: ")

	var response string
	fmt.Scanln(&response)

	if response != "y" && response != "Y" && response != "yes" {
		fmt.Println("❌ Operación cancelada")
		return
	}

	cfg := config.Default()
	if err := cfg.Save(); err != nil {
		fmt.Fprintf(os.Stderr, "❌ Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("✅ Configuración restablecida a valores por defecto")
}
