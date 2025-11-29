package new

import (
	"fmt"
	"os"

	"github.com/okalexiiis/dwrk/internal/project"
	"github.com/okalexiiis/dwrk/pkg/utils"
	"github.com/spf13/cobra"
)

var PROJECTS_DIR = "~/Projects/"

// flags
var (
	git bool
)

var NewCmd = &cobra.Command{
	Use:   "new <nombre>",
	Short: "Crea un nuevo proyecto local",
	Long:  "Crea un nuevo proyecto local en el directorio configurado, opcionalmente instancia un repositorio de git.",
	Args:  cobra.ExactArgs(1), // Requiere exactamente un argumento (el nombre)
	Run:   runNew,
}

func init() {
	NewCmd.Flags().BoolVarP(&git, "git", "g", false, "Inicializar repositorio git")
}

func runNew(cmd *cobra.Command, args []string) {
	projectName := args[0]
	projectsDir := utils.ExpandPath(PROJECTS_DIR)

	// Crear manager
	manager := project.NewManager(projectsDir)

	// Intentar crear el proyecto
	createdProject, err := manager.Create(projectName, project.CreateOptions{
		InitGit: git,
	})

	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Error: %v\n", err)
		os.Exit(1)
	}

	// Mensaje de éxito
	fmt.Printf("✅ Proyecto creado exitosamente: %s\n", createdProject.Name)
	fmt.Printf("📁 Ubicación: %s\n", createdProject.Path)

	if createdProject.IsGit {
		fmt.Println("🔗 Repositorio git inicializado")
	}

	fmt.Printf("\n💡 Para abrir el proyecto:\n")
	fmt.Printf("   dwrk open %s\n", projectName)
}
