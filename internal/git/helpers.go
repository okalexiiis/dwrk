package git

import (
	"fmt"
	"strings"
)

// extractRepoNameFromURL extrae el nombre del repo de una URL
func extractRepoNameFromURL(url string) string {
	// Eliminar .git del final
	url = strings.TrimSuffix(url, ".git")

	// Caso SSH: git@github.com:user/repo
	if strings.Contains(url, ":") {
		parts := strings.Split(url, ":")
		if len(parts) >= 2 {
			url = parts[1]
		}
	}

	// Extraer última parte del path
	parts := strings.Split(url, "/")
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}

	return "cloned-repo"
}

// parseGitError interpreta errores comunes de git
func parseGitError(err error) error {
	errStr := err.Error()

	if strings.Contains(errStr, "Permission denied") {
		return fmt.Errorf("permiso denegado: verifica tu configuración SSH")
	}

	if strings.Contains(errStr, "Repository not found") {
		return fmt.Errorf("repositorio no encontrado: verifica el nombre y permisos")
	}

	if strings.Contains(errStr, "already exists") {
		return fmt.Errorf("el directorio destino ya existe")
	}

	return err
}
