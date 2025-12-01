package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

const (
	// ConfigDirName defines the directory where the configuration file is stored.
	ConfigDirName = ".config/dwrk"

	// ConfigFileName is the filename of the main configuration file.
	ConfigFileName = "config.yaml"
)

// Config represents the application's configuration structure.
type Config struct {
	ProjectsDir    string `yaml:"projects_dir"`    // Local directory where projects are stored.
	TemplatesDir   string `yaml:"templates_dir"`   // Local directiry where templates are stored.
	DefaultEditor  string `yaml:"default_editor"`  // Preferred editor (auto, code, nvim, vim, etc.)
	GitHubUsername string `yaml:"github_username"` // Associated GitHub username.
	UseSSH         bool   `yaml:"use_ssh"`         // Controls whether GitHub operations use SSH.
}

// Default returns a new Config populated with default values.
func Default() *Config {
	homeDir, _ := os.UserHomeDir()
	return &Config{
		ProjectsDir:    filepath.Join(homeDir, "Projects"),
		TemplatesDir:   filepath.Join(homeDir, ".config/dwrk/templates"),
		DefaultEditor:  "auto",
		GitHubUsername: "username",
		UseSSH:         true,
	}
}

// Load loads the configuration from disk.
// If the configuration file does not exist, a new one is created with default values.
func Load() (*Config, error) {
	configPath := GetConfigPath()

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		cfg := Default()

		if err := cfg.Save(); err != nil {
			return nil, fmt.Errorf("failed to create default configuration: %w", err)
		}
		return cfg, nil
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read configuration: %w", err)
	}

	cfg := &Config{}
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("failed to parse configuration: %w", err)
	}

	return cfg, nil
}

// Save writes the current configuration to disk,
// creating the configuration directory if needed.
func (c *Config) Save() error {
	configPath := GetConfigPath()
	configDir := filepath.Dir(configPath)

	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create configuration directory: %w", err)
	}

	data, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Errorf("failed to serialize configuration: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write configuration file: %w", err)
	}

	return nil
}

// Set updates a configuration field by key and saves the result.
// Expected keys: projects_dir, editor, github_username, use_ssh.
func (c *Config) Set(key, value string) error {
	switch key {
	case "projects_dir":
		// Expand '~/' to the user's home directory.
		if len(value) >= 2 && value[:2] == "~/" {
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
		return fmt.Errorf("invalid configuration key: %s", key)
	}

	return c.Save()
}

// Get retrieves the value of a configuration field by key.
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
		return "", fmt.Errorf("invalid configuration key: %s", key)
	}
}

// GetConfigPath returns the absolute path of the configuration file.
func GetConfigPath() string {
	homeDir, _ := os.UserHomeDir()
	return filepath.Join(homeDir, ConfigDirName, ConfigFileName)
}

// Exists returns true if the configuration file exists.
func Exists() bool {
	configPath := GetConfigPath()
	_, err := os.Stat(configPath)
	return err == nil
}
