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

var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "Lista todos los proyectos locales",
	Long:  `Lista todos los proyectos en el directorio de proyectos configurado.`,
	Run:   runList,
}

func init() {
	ListCmd.Flags().BoolVarP(&showHidden, "all", "a", false, "Mostrar carpetas ocultas")
	ListCmd.Flags().StringVarP(&filterName, "filter", "f", "", "Filtrar proyectos por nombre")
}

func runList(cmd *cobra.Command, args []string) {
	// Cargar configuración
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Error cargando configuración: %v\n", err)
		os.Exit(1)
	}

	// Crear manager con el directorio de la config
	manager := project.NewManager(cfg.ProjectsDir)

	// Listar proyectos usando el manager
	projects, err := manager.List(project.ListOptions{
		ShowHidden: showHidden,
		Filter:     filterName,
	})

	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Error al listar proyectos: %v\n", err)
		os.Exit(1)
	}

	if len(projects) == 0 {
		fmt.Printf("No se encontraron proyectos en: %s\n", cfg.ProjectsDir)
		fmt.Println("\n💡 Crea un proyecto nuevo:")
		fmt.Println("   dwrk new mi-proyecto")
		return
	}

	fmt.Printf("📁 Proyectos en %s:\n\n", cfg.ProjectsDir)
	for i, proj := range projects {
		// Indicador de git
		gitIndicator := ""
		if proj.IsGit {
			gitIndicator = " 🔗"
		}

		fmt.Printf("  %d. %s%s\n", i+1, proj.Name, gitIndicator)
	}

	fmt.Printf("\nTotal: %d proyecto(s)\n", len(projects))
}
