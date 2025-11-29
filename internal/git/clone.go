package git

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

type Cloner struct{}

func NewCloner() *Cloner {
	return &Cloner{}
}

// Clone clona un repositorio en el directorio especificado
func (c *Cloner) Clone(repoURL, targetDir string) (string, error) {
	// Verificar que git esté instalado
	if !c.IsGitInstalled() {
		return "", fmt.Errorf("git no está instalado en el sistema")
	}

	// Crear directorio destino si no existe
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return "", fmt.Errorf("error creando directorio destino: %w", err)
	}

	// Ejecutar git clone
	cmd := exec.Command("git", "clone", repoURL)
	cmd.Dir = targetDir

	// Capturar salida para mostrar progreso
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return "", parseGitError(err)
	}

	// Determinar el nombre del directorio clonado
	repoName := extractRepoNameFromURL(repoURL)
	clonedPath := filepath.Join(targetDir, repoName)

	return clonedPath, nil
}

// CloneWithDepth clona con profundidad limitada (shallow clone)
func (c *Cloner) CloneWithDepth(repoURL, targetDir string, depth int) (string, error) {
	if !c.IsGitInstalled() {
		return "", fmt.Errorf("git no está instalado en el sistema")
	}

	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return "", fmt.Errorf("error creando directorio destino: %w", err)
	}

	cmd := exec.Command("git", "clone", "--depth", fmt.Sprintf("%d", depth), repoURL)
	cmd.Dir = targetDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return "", parseGitError(err)
	}

	repoName := extractRepoNameFromURL(repoURL)
	clonedPath := filepath.Join(targetDir, repoName)

	return clonedPath, nil
}

// CloneBranch clona una rama específica
func (c *Cloner) CloneBranch(repoURL, targetDir, branch string) (string, error) {
	if !c.IsGitInstalled() {
		return "", fmt.Errorf("git no está instalado en el sistema")
	}

	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return "", fmt.Errorf("error creando directorio destino: %w", err)
	}

	cmd := exec.Command("git", "clone", "-b", branch, repoURL)
	cmd.Dir = targetDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return "", parseGitError(err)
	}

	repoName := extractRepoNameFromURL(repoURL)
	clonedPath := filepath.Join(targetDir, repoName)

	return clonedPath, nil
}

// IsGitInstalled verifica si git está instalado
func (c *Cloner) IsGitInstalled() bool {
	_, err := exec.LookPath("git")
	return err == nil
}
