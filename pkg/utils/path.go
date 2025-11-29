package utils

import (
	"os"
	"path/filepath"
	"strings"
)

// ExpandPath expande ~ al directorio home del usuario
func ExpandPath(path string) string {
	if strings.HasPrefix(path, "~/") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return path
		}
		return filepath.Join(homeDir, path[2:])
	}
	return path
}
