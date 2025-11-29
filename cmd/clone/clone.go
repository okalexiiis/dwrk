// ============================================
// cmd/clone/clone.go
// ============================================
package clone

import (
	"fmt"
	"os"
	"strings"

	"github.com/okalexiiis/dwrk/internal/git"
	"github.com/okalexiiis/dwrk/internal/project"
	"github.com/okalexiiis/dwrk/pkg/utils"
	"github.com/spf13/cobra"
)

var PROJECTS_DIR = "~/Projects/"
var DEFAULT_USERNAME = "okalexiiis"

// flags
var (
	username string
	url      string
	useHTTPS bool
	destDir  string
)

var CloneCmd = &cobra.Command{
	Use:   "clone <repo>",
	Short: "Clona un repositorio de GitHub",
	Long: `Clona un repositorio de GitHub en el directorio de proyectos.

Ejemplos:
  proj clone my-repo                    # Clona tu repo con SSH
  proj clone my-repo --https            # Clona tu repo con HTTPS
  proj clone my-repo -u otro-usuario    # Clona repo de otro usuario
  proj clone --url https://github.com/user/repo.git  # Clona URL específica
  proj clone my-repo --dir ~/Otros      # Clona en directorio diferente`,
	Args: cobra.MaximumNArgs(1),
	Run:  runClone,
}

func init() {
	CloneCmd.Flags().StringVarP(&username, "user", "u", DEFAULT_USERNAME, "Usuario de GitHub")
	CloneCmd.Flags().StringVar(&url, "url", "", "URL completa del repositorio (ignora otras flags)")
	CloneCmd.Flags().BoolVar(&useHTTPS, "https", false, "Usar HTTPS en lugar de SSH")
	CloneCmd.Flags().StringVar(&destDir, "dir", "", "Directorio destino (por defecto: PROJECTS_DIR)")
}

func runClone(cmd *cobra.Command, args []string) {
	var repoURL string
	var repoName string

	// Determinar la URL del repositorio
	if url != "" {
		// Caso 1: URL completa proporcionada
		repoURL = url
		repoName = extractRepoName(url)
		fmt.Printf("🔗 Clonando desde URL: %s\n", repoURL)
	} else {
		// Caso 2: Construir URL desde nombre de repo
		if len(args) == 0 {
			fmt.Fprintln(os.Stderr, "❌ Error: debes proporcionar un nombre de repositorio o usar --url")
			fmt.Println("\n💡 Ejemplos:")
			fmt.Println("   proj clone my-repo")
			fmt.Println("   proj clone --url https://github.com/user/repo.git")
			os.Exit(1)
		}

		repoName = args[0]
		repoURL = buildRepoURL(username, repoName, useHTTPS)

		protocol := "SSH"
		if useHTTPS {
			protocol = "HTTPS"
		}
		fmt.Printf("🔗 Clonando %s/%s (%s)...\n", username, repoName, protocol)
	}

	// Determinar directorio destino
	targetDir := PROJECTS_DIR
	if destDir != "" {
		targetDir = destDir
	}
	targetPath := utils.ExpandPath(targetDir)

	// Verificar que no exista ya el proyecto
	manager := project.NewManager(targetPath)
	if manager.Exists(repoName) {
		fmt.Fprintf(os.Stderr, "❌ Error: ya existe un proyecto con el nombre '%s'\n", repoName)
		fmt.Printf("📁 Ubicación: %s/%s\n", targetPath, repoName)
		fmt.Println("\n💡 Usa --dir para clonar en otra ubicación")
		os.Exit(1)
	}

	// Clonar el repositorio
	cloner := git.NewCloner()
	clonedPath, err := cloner.Clone(repoURL, targetPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Error al clonar repositorio: %v\n", err)

		// Sugerencias según el tipo de error
		if strings.Contains(err.Error(), "Permission denied") {
			fmt.Println("\n💡 Si usas SSH, verifica que tu clave esté configurada:")
			fmt.Println("   ssh -T git@github.com")
			fmt.Println("   O intenta con HTTPS: --https")
		} else if strings.Contains(err.Error(), "Repository not found") {
			fmt.Println("\n💡 Verifica que:")
			fmt.Println("   - El repositorio existe")
			fmt.Println("   - Tienes permisos de acceso")
			fmt.Println("   - El nombre de usuario es correcto")
		}

		os.Exit(1)
	}

	// Éxito
	fmt.Printf("✅ Repositorio clonado exitosamente\n")
	fmt.Printf("📁 Ubicación: %s\n", clonedPath)
	fmt.Printf("\n💡 Para abrir el proyecto:\n")
	fmt.Printf("   proj open %s\n", repoName)
}

// buildRepoURL construye la URL del repositorio según el protocolo
func buildRepoURL(user, repo string, https bool) string {
	if https {
		return fmt.Sprintf("https://github.com/%s/%s.git", user, repo)
	}
	return fmt.Sprintf("git@github.com:%s/%s.git", user, repo)
}

// extractRepoName extrae el nombre del repositorio de una URL
func extractRepoName(url string) string {
	// Eliminar .git del final si existe
	url = strings.TrimSuffix(url, ".git")

	// Extraer última parte de la URL
	parts := strings.Split(url, "/")
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}

	return "cloned-repo"
}
