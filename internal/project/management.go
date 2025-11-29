package project

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type Manager struct {
	baseDir string
}

type ListOptions struct {
	ShowHidden bool
	Filter     string
	Search     string
}

type CreateOptions struct {
	InitGit bool
}

type Project struct {
	Name    string
	Path    string
	IsGit   bool
	LastMod time.Time
}

func NewManager(baseDir string) *Manager {
	return &Manager{baseDir: baseDir}
}

// List lista todos los proyectos
func (m *Manager) List(opts ListOptions) ([]Project, error) {
	entries, err := os.ReadDir(m.baseDir)
	if err != nil {
		return nil, err
	}

	var projects []Project
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		name := entry.Name()

		if !opts.ShowHidden && strings.HasPrefix(name, ".") {
			continue
		}

		if opts.Filter != "" && !strings.Contains(strings.ToLower(name), strings.ToLower(opts.Filter)) {
			continue
		}

		projectPath := filepath.Join(m.baseDir, name)
		info, _ := entry.Info()

		projects = append(projects, Project{
			Name:    name,
			Path:    projectPath,
			IsGit:   isGitRepo(projectPath),
			LastMod: info.ModTime(),
		})
	}

	return projects, nil
}

// Create crea un nuevo proyecto
func (m *Manager) Create(name string, opts CreateOptions) (*Project, error) {
	// 1. Validar nombre del proyecto
	if err := validateProjectName(name); err != nil {
		return nil, err
	}

	// 2. Verificar si ya existe
	if m.Exists(name) {
		return nil, fmt.Errorf("ya existe un proyecto con el nombre '%s'", name)
	}

	// 3. Crear el directorio del proyecto
	projectPath := filepath.Join(m.baseDir, name)
	if err := os.MkdirAll(projectPath, 0755); err != nil {
		return nil, fmt.Errorf("error al crear directorio: %w", err)
	}

	// 4. Inicializar git si se solicita
	isGit := false
	if opts.InitGit {
		if err := initGitRepo(projectPath, name); err != nil {
			// Intentamos limpiar si falla
			os.RemoveAll(projectPath)
			return nil, fmt.Errorf("error al inicializar git: %w", err)
		}
		isGit = true
	}

	// 5. Retornar el proyecto creado
	return &Project{
		Name:    name,
		Path:    projectPath,
		IsGit:   isGit,
		LastMod: time.Now(),
	}, nil
}

// Exists verifica si un proyecto existe
func (m *Manager) Exists(name string) bool {
	projectPath := filepath.Join(m.baseDir, name)
	info, err := os.Stat(projectPath)
	if err != nil {
		return false
	}
	return info.IsDir()
}

// Get obtiene información de un proyecto específico
func (m *Manager) Get(name string) (*Project, error) {
	if !m.Exists(name) {
		return nil, fmt.Errorf("proyecto '%s' no encontrado", name)
	}

	projectPath := filepath.Join(m.baseDir, name)
	info, err := os.Stat(projectPath)
	if err != nil {
		return nil, err
	}

	return &Project{
		Name:    name,
		Path:    projectPath,
		IsGit:   isGitRepo(projectPath),
		LastMod: info.ModTime(),
	}, nil
}

// isGitRepo verifica si un directorio es un repositorio git
func isGitRepo(path string) bool {
	gitPath := filepath.Join(path, ".git")
	info, err := os.Stat(gitPath)
	return err == nil && info.IsDir()
}

// validateProjectName valida que el nombre sea válido
func validateProjectName(name string) error {
	if name == "" {
		return fmt.Errorf("el nombre del proyecto no puede estar vacío")
	}

	// Caracteres no permitidos en nombres de carpetas
	invalidChars := []string{"/", "\\", ":", "*", "?", "\"", "<", ">", "|"}
	for _, char := range invalidChars {
		if strings.Contains(name, char) {
			return fmt.Errorf("el nombre contiene caracteres no permitidos: %s", char)
		}
	}

	// No permitir nombres que empiecen con punto (carpetas ocultas)
	if strings.HasPrefix(name, ".") {
		return fmt.Errorf("el nombre no puede empezar con punto")
	}

	return nil
}

// initGitRepo inicializa un repositorio git y crea README
func initGitRepo(projectPath, projectName string) error {
	// Ejecutar git init
	cmd := exec.Command("git", "init")
	cmd.Dir = projectPath
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error ejecutando 'git init': %w", err)
	}

	// Crear README.md
	readmePath := filepath.Join(projectPath, "README.md")
	readmeContent := fmt.Sprintf("# %s\n\nProyecto creado con Project Manager CLI.\n", projectName)

	if err := os.WriteFile(readmePath, []byte(readmeContent), 0644); err != nil {
		return fmt.Errorf("error creando README.md: %w", err)
	}

	// Hacer commit inicial (opcional pero recomendado)
	cmd = exec.Command("git", "add", "README.md")
	cmd.Dir = projectPath
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error en 'git add': %w", err)
	}

	cmd = exec.Command("git", "commit", "-m", "Initial commit")
	cmd.Dir = projectPath
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error en 'git commit': %w", err)
	}

	return nil
}
