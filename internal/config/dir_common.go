//go:build !windows

package config

import (
	"fmt"
	"os"
	"path/filepath"
)

const DirectoryName = ".githoard"

func DirectoryPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("cannot resolve user home: %w", err)
	}

	location := filepath.Join(home, DirectoryName)

	stat, err := os.Stat(location)
	if err != nil && os.IsNotExist(err) {
		if err := os.MkdirAll(location, 0772); err != nil {
			return "", fmt.Errorf("failed to create directory: %w", err)
		}
	} else if err != nil {
		return "", fmt.Errorf("failed to stat directory: %w", err)
	} else if !stat.IsDir() {
		return "", fmt.Errorf("file collision: %s is a regular file", location)
	}

	return home, nil
}
