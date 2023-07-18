package main

import (
	"os"
	"path"
	"path/filepath"
	"strings"
)

func applyCommands(commands map[string]string, result transformed) (transformed, error) {
	if commands["extends"] != "" {
		currentDirectory := strings.Join(strings.Split(result.Path, "/")[:len(strings.Split(result.Path, "/"))-1], "/")
		templatePath, err := filepath.Abs(path.Join(currentDirectory, commands["extends"]))
		if err != nil {
			return transformed{}, err
		}
		by, err := os.ReadFile(templatePath)
		if err != nil {
			return transformed{}, err
		}
		result.Extends = string(by)
	}

	return result, nil
}
