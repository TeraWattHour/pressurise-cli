package main

import (
	"os"
	"path"
	"path/filepath"
	"strings"
)

func applyCommands(commands map[string][]string, result transformed) (transformed, error) {
	currentDirectory := strings.Join(strings.Split(result.Path, "/")[:len(strings.Split(result.Path, "/"))-1], "/")

	if len(commands["extends"]) != 0 {
		templatePath, err := filepath.Abs(path.Join(currentDirectory, commands["extends"][0]))
		if err != nil {
			return transformed{}, err
		}
		by, err := os.ReadFile(templatePath)
		if err != nil {
			return transformed{}, err
		}
		result.Extends = string(by)
	}

	if len(commands["component"]) != 0 {
		for _, component := range commands["component"] {
			componentPath, err := filepath.Abs(path.Join(currentDirectory, component))
			if err != nil {
				return transformed{}, err
			}
			by, err := os.ReadFile(componentPath)
			if err != nil {
				return transformed{}, err
			}

			result.Components = append(result.Components, string(by))
		}

	}

	return result, nil
}
