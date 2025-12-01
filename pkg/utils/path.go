package utils

import (
	"fmt"
	"io"
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

// CopyDir recursively copies the contents of the source directory (og) to the
// destination directory (to).
// It ensures that all file permissions and directory structure are preserved.
func CopyDir(og, to string) error {
	// 1. Get info about the source directory
	sourceInfo, err := os.Stat(og)
	if err != nil {
		return fmt.Errorf("stat source: %w", err)
	}

	// 2. Create the destination directory with the source's permissions
	if err := os.MkdirAll(to, sourceInfo.Mode()); err != nil {
		return fmt.Errorf("mkdir destination: %w", err)
	}

	// 3. Read the contents of the source directory
	entries, err := os.ReadDir(og)
	if err != nil {
		return fmt.Errorf("readdir source: %w", err)
	}

	// 4. Iterate over all entries and copy them
	for _, entry := range entries {
		sourcePath := filepath.Join(og, entry.Name())
		destPath := filepath.Join(to, entry.Name())

		// If it's a directory, recursively call CopyDir
		if entry.IsDir() {
			if err := CopyDir(sourcePath, destPath); err != nil {
				return err
			}
			continue
		}

		// If it's a file, copy the file content
		if err := copyFile(sourcePath, destPath); err != nil {
			return err
		}
	}

	return nil
}

// copyFile copies a single file from sourcePath to destPath.
func copyFile(sourcePath, destPath string) error {
	source, err := os.Open(sourcePath)
	if err != nil {
		return fmt.Errorf("open source file: %w", err)
	}
	defer source.Close()

	// Get file info to preserve permissions
	sourceStat, err := source.Stat()
	if err != nil {
		return fmt.Errorf("stat source file: %w", err)
	}

	// Create destination file with the same permissions
	dest, err := os.OpenFile(destPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, sourceStat.Mode())
	if err != nil {
		return fmt.Errorf("create dest file: %w", err)
	}
	defer dest.Close()

	// Copy content
	if _, err := io.Copy(dest, source); err != nil {
		return fmt.Errorf("copy content: %w", err)
	}

	return nil
}
