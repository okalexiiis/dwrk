package clone

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/okalexiiis/dwrk/internal/config"
	"github.com/okalexiiis/dwrk/internal/git"
	"github.com/okalexiiis/dwrk/internal/project"
	"github.com/okalexiiis/dwrk/pkg/utils"
	"github.com/spf13/cobra"
)

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
	Args:  cobra.MaximumNArgs(1),
	Run:   runClone,
}

func init() {
	CloneCmd.Flags().StringVarP(&username, "user", "u", "", "Usuario de GitHub (por defecto: de config)")
	CloneCmd.Flags().StringVar(&url, "url", "", "URL completa del repositorio")
	CloneCmd.Flags().BoolVar(&useHTTPS, "https", false, "Usar HTTPS en lugar de SSH")
	CloneCmd.Flags().StringVar(&destDir, "dir", "", "Directorio destino")
}

func runClone(cmd *cobra.Command, args []string) {
	// Cargar configuración
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Error cargando configuración: %v\n", err)
		os.Exit(1)
	}

	// Usar username de config si no se especificó
	if username == "" {
		username = cfg.GitHubUsername
	}

	// Usar preferencia de SSH/HTTPS de config si no se especificó
	if !cmd.Flags().Changed("https") {
		useHTTPS = !cfg.UseSSH
	}

	var repoURL string
	var repoName string

	if url != "" {
		repoURL = url
		repoName = utils.ExtractRepoNameFromURL(url)
		fmt.Printf("🔗 Clonando desde URL: %s\n", repoURL)
	} else {
		if len(args) == 0 {
			fmt.Fprintln(os.Stderr, "❌ Error: debes proporcionar un nombre de repositorio o usar --url")
			os.Exit(1)
		}

		repoName = args[0]

		// Usar PROJECTS_DIR de config
		projectsPath := cfg.ProjectsDir
		manager := project.NewManager(projectsPath)

		if destDir == "" {
			if manager.Exists(repoName) {
				destDir = filepath.Join(cfg.ProjectsDir, repoName)
				fmt.Printf("📁 Encontrado proyecto local '%s'\n", repoName)
				fmt.Printf("📥 Clonando dentro de: %s\n", destDir)
			} else {
				fmt.Printf("⚠️  No existe un proyecto local llamado '%s'\n", repoName)
				fmt.Print("¿Deseas clonar en PROJECTS_DIR y crear el directorio? [Y/n]: ")

				var response string
				fmt.Scanln(&response)

				response = strings.ToLower(strings.TrimSpace(response))
				if response == "n" || response == "no" {
					fmt.Println("❌ Operación cancelada")
					os.Exit(0)
				}

				destDir = cfg.ProjectsDir
			}
		}

		repoURL = utils.BuildRepoURL(username, repoName, useHTTPS)

		protocol := "SSH"
		if useHTTPS {
			protocol = "HTTPS"
		}
		fmt.Printf("🔗 Clonando %s/%s (%s)...\n", username, repoName, protocol)
	}

	targetPath := utils.ExpandPath(destDir)

	if err := os.MkdirAll(targetPath, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "❌ Error: no se pudo crear el directorio destino: %v\n", err)
		os.Exit(1)
	}

	cloner := git.NewCloner()
	clonedPath, err := cloner.Clone(repoURL, targetPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Error al clonar repositorio: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✅ Repositorio clonado exitosamente\n")
	fmt.Printf("📁 Ubicación: %s\n", clonedPath)
	fmt.Printf("\n💡 Para abrir el proyecto:\n")
	fmt.Printf("   dwrk open %s\n", repoName)
}
