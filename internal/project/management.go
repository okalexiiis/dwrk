package project

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// Manager handles project discovery, creation, and metadata retrieval.
type Manager struct {
	baseDir string
}

// ListOptions defines filtering options for the List method.
type ListOptions struct {
	ShowHidden bool   // Include directories starting with a dot
	Filter     string // Substring filter applied to project names
	Search     string // Future extension for fuzzy search or advanced matching
}

// CreateOptions defines optional behaviors when creating a project.
type CreateOptions struct {
	InitGit bool // Initialize a Git repository after creating the folder
}

// Project describes a project discovered or created by the Manager.
type Project struct {
	Name    string
	Path    string
	IsGit   bool
	LastMod time.Time
}

// NewManager creates a new Manager using the provided base directory.
func NewManager(baseDir string) *Manager {
	return &Manager{baseDir: baseDir}
}

// List returns all projects under the base directory, applying the given filters.
//
// A project is any non-hidden directory unless ShowHidden is enabled.
// Git repositories are automatically detected.
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

		if opts.Filter != "" &&
			!strings.Contains(strings.ToLower(name), strings.ToLower(opts.Filter)) {
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

// Create creates a new project directory and optionally initializes a Git repository.
//
// The name is validated to avoid invalid or unsafe folder names.
// If Git initialization fails, the created directory is removed to maintain consistency.
func (m *Manager) Create(name string, opts CreateOptions) (*Project, error) {
	if err := validateProjectName(name); err != nil {
		return nil, err
	}

	if m.Exists(name) {
		return nil, fmt.Errorf("a project named '%s' already exists", name)
	}

	projectPath := filepath.Join(m.baseDir, name)
	if err := os.MkdirAll(projectPath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create directory: %w", err)
	}

	isGit := false
	if opts.InitGit {
		if err := initGitRepo(projectPath, name); err != nil {
			os.RemoveAll(projectPath)
			return nil, fmt.Errorf("failed to initialize git: %w", err)
		}
		isGit = true
	}

	return &Project{
		Name:    name,
		Path:    projectPath,
		IsGit:   isGit,
		LastMod: time.Now(),
	}, nil
}

// Exists checks whether a project with the given name exists.
func (m *Manager) Exists(name string) bool {
	projectPath := filepath.Join(m.baseDir, name)
	info, err := os.Stat(projectPath)
	if err != nil {
		return false
	}
	return info.IsDir()
}

// Get retrieves metadata about a specific project by name.
func (m *Manager) Get(name string) (*Project, error) {
	if !m.Exists(name) {
		return nil, fmt.Errorf("project '%s' not found", name)
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

// isGitRepo returns true if the given path contains a .git folder.
func isGitRepo(path string) bool {
	gitPath := filepath.Join(path, ".git")
	info, err := os.Stat(gitPath)
	return err == nil && info.IsDir()
}

// validateProjectName ensures the project name is safe and valid for use as a directory name.
func validateProjectName(name string) error {
	if name == "" {
		return fmt.Errorf("project name cannot be empty")
	}

	invalidChars := []string{"/", "\\", ":", "*", "?", "\"", "<", ">", "|"}
	for _, char := range invalidChars {
		if strings.Contains(name, char) {
			return fmt.Errorf("project name contains invalid character: %s", char)
		}
	}

	if strings.HasPrefix(name, ".") {
		return fmt.Errorf("project name cannot start with a dot")
	}

	return nil
}

// initGitRepo initializes a Git repository and generates a README.md file.
//
// It runs `git init`, creates a default README, adds it to the repository,
// and performs an initial commit.
func initGitRepo(projectPath, projectName string) error {
	cmd := exec.Command("git", "init")
	cmd.Dir = projectPath
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed running 'git init': %w", err)
	}

	// Generate README.md
	readmePath := filepath.Join(projectPath, "README.md")
	readmeContent := fmt.Sprintf("# %s\n\nProject initialized using Project Manager CLI.\n", projectName)

	if err := os.WriteFile(readmePath, []byte(readmeContent), 0644); err != nil {
		return fmt.Errorf("failed creating README.md: %w", err)
	}

	cmd = exec.Command("git", "add", "README.md")
	cmd.Dir = projectPath
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed running 'git add': %w", err)
	}

	cmd = exec.Command("git", "commit", "-m", "Initial commit")
	cmd.Dir = projectPath
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed running 'git commit': %w", err)
	}

	return nil
}
