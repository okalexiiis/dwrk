package list

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var PROJECT_DIR = "~/Projects/"

// Flags
var (
	showHidden bool
	filterName string
)

var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "Lista todos los proyectos locales",
	Long:  `Lista todos los proyectos en el directorio de proyectos configurado.`,
	Run: func(cmd *cobra.Command, args []string) {
		projects, err := getProjects()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error al listar proyectos: %v\n", err)
			os.Exit(1)
		}

		if len(projects) == 0 {
			fmt.Println("No se encontraron proyectos en:", expandPath(PROJECT_DIR))
			return
		}

		fmt.Printf("📁 Proyectos en %s:\n\n", expandPath(PROJECT_DIR))
		for i, project := range projects {
			// Verificar si es un repo git
			gitIndicator := ""
			if isGitRepo(filepath.Join(expandPath(PROJECT_DIR), project)) {
				gitIndicator = " 🔗"
			}

			fmt.Printf("  %d. %s%s\n", i+1, project, gitIndicator)
		}
		fmt.Printf("\nTotal: %d proyecto(s)\n", len(projects))
	},
}

func init() {
	ListCmd.Flags().BoolVarP(&showHidden, "all", "a", false, "Mostrar carpetas ocultas")
	ListCmd.Flags().StringVarP(&filterName, "filter", "f", "", "Filtrar proyectos por nombre")
}

// getProjects obtiene la lista de directorios en PROJECT_DIR
func getProjects() ([]string, error) {
	projectPath := expandPath(PROJECT_DIR)

	// Verificar si el directorio existe
	if _, err := os.Stat(projectPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("el directorio de proyectos no existe: %s", projectPath)
	}

	entries, err := os.ReadDir(projectPath)
	if err != nil {
		return nil, err
	}

	var projects []string
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		name := entry.Name()

		// Filtrar carpetas ocultas si no se especifica --all
		if !showHidden && strings.HasPrefix(name, ".") {
			continue
		}

		// Aplicar filtro de nombre si existe
		if filterName != "" && !strings.Contains(strings.ToLower(name), strings.ToLower(filterName)) {
			continue
		}

		projects = append(projects, name)
	}

	return projects, nil
}

// expandPath expande ~ al directorio home del usuario
func expandPath(path string) string {
	if strings.HasPrefix(path, "~/") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return path
		}
		return filepath.Join(homeDir, path[2:])
	}
	return path
}

// isGitRepo verifica si un directorio es un repositorio git
func isGitRepo(path string) bool {
	gitPath := filepath.Join(path, ".git")
	info, err := os.Stat(gitPath)
	if err != nil {
		return false
	}
	return info.IsDir()
}
