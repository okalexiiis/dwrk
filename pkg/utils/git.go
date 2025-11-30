package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// extractRepoName extrae el nombre del repositorio de una URL
func ExtractRepoName(url string) string {
	// Eliminar .git del final si existe
	url = strings.TrimSuffix(url, ".git")

	// Extraer última parte de la URL
	parts := strings.Split(url, "/")
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}

	return "cloned-repo"
}

// buildRepoURL construye la URL del repositorio según el protocolo
func BuildRepoURL(user, repo string, https bool) string {
	if https {
		return fmt.Sprintf("https://github.com/%s/%s.git", user, repo)
	}
	return fmt.Sprintf("git@github.com:%s/%s.git", user, repo)
}

// isGitRepo verifica si un directorio es un repositorio git
func IsGitRepo(path string) bool {
	gitPath := filepath.Join(path, ".git")
	info, err := os.Stat(gitPath)
	if err != nil {
		return false
	}
	return info.IsDir()
}
