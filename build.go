package main

import (
	"bytes"
	"errors"
	"fmt"
	"go/ast"
	"os"
	"path/filepath"
	"text/template"

	"golang.org/x/exp/slices"
)

var fileTemplateRaw = `// Don't edit this file!
// This file is generated, if you run the generator again
// all changes made to this file will be lost.

package main

{{ .Imports }}

var pressHandlers = map[string]func() http.HandlerFunc {
    {{ .Handlers }}
}`

var fileTemplate = template.Must(template.New("fileTemplate").Parse(fileTemplateRaw))

var fileImports = []*ast.ImportSpec{{
	Name: &ast.Ident{Name: "html"},
	Path: &ast.BasicLit{Value: "\"html/template\""},
}, {
	Name: &ast.Ident{Name: "http"},
	Path: &ast.BasicLit{Value: "\"net/http\""},
}}

func buildProject(directory string) error {
	err := verifyProjectDirectory(directory)
	if err != nil {
		return err
	}

	appDirectory := filepath.Join(directory, "app")

	foundPages, err := scanForPages(appDirectory, []string{})
	if err != nil {
		return err
	}

	transformed := []transformed{}

	for _, page := range foundPages {
		result := transformPageFile(page, appDirectory)
		transformed = append(transformed, result)
	}

	imports := fileImports
	handlers := ""
	usedHandlers := []string{}

	for _, page := range transformed {
		handlers += page.Handler
		if slices.Contains(usedHandlers, page.RouterPath) {
			return fmt.Errorf("conflicting paths, `%s` path is used at more than one place", page.RouterPath)
		}
		usedHandlers = append(usedHandlers, page.RouterPath)

	pageLoop:
		for _, spec := range page.Imports {
			if spec == nil {
				continue
			}
			for _, declared := range imports {
				if spec.Name != nil && declared.Name == spec.Name {
					if declared.Path.Value != spec.Path.Value {
						return fmt.Errorf("conflicting import, import alias `%s` is already being used for package `%s` and is duplicated by `%s`", declared.Name, declared.Path.Value, spec.Path.Value)
					}
					break pageLoop
				}
			}
			imports = append(imports, spec)
		}
	}

	var tpl bytes.Buffer
	if err := fileTemplate.Execute(&tpl, map[string]interface{}{
		"Handlers": handlers,
		"Imports":  makeImports(imports),
	}); err != nil {
		panic(err)
	}

	outputFile := filepath.Join(directory, "pressrelease.go")
	if err := os.Remove(outputFile); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return err
		}
	}

	file, err := os.Create(outputFile)
	if err != nil {
		return err
	}

	if _, err := file.WriteString(tpl.String()); err != nil {
		return err
	}

	return nil
}
