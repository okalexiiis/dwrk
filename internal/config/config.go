package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

const (
	ConfigDirName  = ".config/dwrk"
	ConfigFileName = "config.yaml"
)

// Config estructura principal de configuración
type Config struct {
	ProjectsDir    string `yaml:"projects_dir"`
	DefaultEditor  string `yaml:"default_editor"`
	GitHubUsername string `yaml:"github_username"`
	UseSSH         bool   `yaml:"use_ssh"`
}

// Default retorna la configuración por defecto
func Default() *Config {
	homeDir, _ := os.UserHomeDir()
	return &Config{
		ProjectsDir:    filepath.Join(homeDir, "Projects"),
		DefaultEditor:  "auto", // auto, code, nvim, vim, nano, terminal
		GitHubUsername: "username",
		UseSSH:         true,
	}
}

// Load carga la configuración desde el archivo
func Load() (*Config, error) {
	configPath := GetConfigPath()

	// Si no existe el archivo, crear con valores por defecto
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		cfg := Default()
		if err := cfg.Save(); err != nil {
			return nil, fmt.Errorf("error creando configuración por defecto: %w", err)
		}
		return cfg, nil
	}

	// Leer archivo
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("error leyendo configuración: %w", err)
	}

	// Parsear YAML
	cfg := &Config{}
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("error parseando configuración: %w", err)
	}

	return cfg, nil
}

// Save guarda la configuración en el archivo
func (c *Config) Save() error {
	configPath := GetConfigPath()

	// Crear directorio si no existe
	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("error creando directorio de configuración: %w", err)
	}

	// Convertir a YAML
	data, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Errorf("error serializando configuración: %w", err)
	}

	// Escribir archivo
	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("error escribiendo configuración: %w", err)
	}

	return nil
}

// Set actualiza un valor de configuración
func (c *Config) Set(key, value string) error {
	switch key {
	case "projects_dir":
		// Expandir ~ si es necesario
		if value[:2] == "~/" {
			homeDir, _ := os.UserHomeDir()
			value = filepath.Join(homeDir, value[2:])
		}
		c.ProjectsDir = value
	case "editor", "default_editor":
		c.DefaultEditor = value
	case "github_username", "username":
		c.GitHubUsername = value
	case "use_ssh", "ssh":
		c.UseSSH = value == "true" || value == "yes" || value == "1"
	default:
		return fmt.Errorf("clave de configuración no válida: %s", key)
	}

	return c.Save()
}

// Get obtiene un valor de configuración
func (c *Config) Get(key string) (string, error) {
	switch key {
	case "projects_dir":
		return c.ProjectsDir, nil
	case "editor", "default_editor":
		return c.DefaultEditor, nil
	case "github_username", "username":
		return c.GitHubUsername, nil
	case "use_ssh", "ssh":
		if c.UseSSH {
			return "true", nil
		}
		return "false", nil
	default:
		return "", fmt.Errorf("clave de configuración no válida: %s", key)
	}
}

// GetConfigPath retorna la ruta del archivo de configuración
func GetConfigPath() string {
	homeDir, _ := os.UserHomeDir()
	return filepath.Join(homeDir, ConfigDirName, ConfigFileName)
}

// Exists verifica si existe el archivo de configuración
func Exists() bool {
	configPath := GetConfigPath()
	_, err := os.Stat(configPath)
	return err == nil
}
