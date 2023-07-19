package main

import (
	"errors"
	"os"
	"path"
)

func verifyProjectDirectory(directory string) error {
	stat, err := os.Stat(path.Join(directory, "app"))
	if err != nil {
		if errors.Is(err, os.ErrNotExist) || stat.IsDir() {
			return errors.New("required directory `app/` doesn't exist in the project directory")
		}
		return err
	}

	return nil
}

func scanPagesDirectory(directory string, found []string) ([]string, error) {
	entries, err := os.ReadDir(directory)
	if err != nil {
		return nil, err
	}
	for _, entry := range entries {
		if entry.IsDir() {
			found, err = scanPagesDirectory(path.Join(directory, entry.Name()), found)
			if err != nil {
				return nil, err
			}
		} else {
			if len(entry.Name()) <= 5 || entry.Name()[len(entry.Name())-5:] != ".html" {
				continue
			}
			p := path.Join(directory, entry.Name())
			found = append(found, p)
		}
	}
	return found, nil
}
