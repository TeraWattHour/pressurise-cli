package main

import (
	"errors"
	"os"
	"path/filepath"
)

func verifyProjectDirectory(directory string) error {
	stat, err := os.Stat(filepath.Join(directory, "app"))
	if err != nil {
		if errors.Is(err, os.ErrNotExist) || stat.IsDir() {
			return errors.New("required directory `app/` doesn't exist in the project directory")
		}
		return err
	}

	return nil
}

func scanForPages(directory string, found []string) ([]string, error) {
	entries, err := os.ReadDir(directory)
	if err != nil {
		return nil, err
	}
	for _, entry := range entries {
		currentPath := filepath.Join(directory, entry.Name())
		if entry.IsDir() {
			found, err = scanForPages(currentPath, found)
			if err != nil {
				return nil, err
			}
		} else {
			if len(entry.Name()) <= 5 || entry.Name()[len(entry.Name())-5:] != ".html" {
				continue
			}
			found = append(found, currentPath)
		}
	}
	return found, nil
}
