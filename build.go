package main

import (
	"bytes"
	"errors"
	"fmt"
	"go/ast"
	"os"
	"path"
	"text/template"
)

var fileTemplate = template.Must(template.ParseFiles("file_template.txt"))
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

	appDirectory := path.Join(directory, "app")

	foundPages, err := scanPagesDirectory(appDirectory, []string{})
	if err != nil {
		return err
	}

	transformed := []transformed{}

	for _, page := range foundPages {
		result := transformPageFile(appDirectory, page)
		transformed = append(transformed, result)
	}

	imports := fileImports
	handlers := ""

	for _, page := range transformed {
		handlers += page.Handler

	pageLoop:
		for _, spec := range page.Imports {
			if spec == nil {
				continue
			}
			for _, declared := range imports {
				if spec.Name != nil && declared.Name == spec.Name {
					if declared.Path.Value != spec.Path.Value {
						return fmt.Errorf("conflicting import, import alias `%s` already used for package `%s` and duplicated by `%s`", declared.Name, declared.Path.Value, spec.Path.Value)
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

	outputFile := path.Join(directory, "pressrelease.go")
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
