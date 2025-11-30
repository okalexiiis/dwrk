package utils

import (
	"os"
	"path/filepath"
	"strings"
)

// Expands the ~ dit to the full path /home/user/
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
