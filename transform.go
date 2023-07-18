package main

import (
	"bytes"
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"strings"
	"text/template"

	"golang.org/x/exp/slices"
)

var handlerTemplate = template.Must(template.ParseFiles("handler_template.txt"))

type transformed struct {
	RouterPath string
	Path       string

	RawCode string

	Imports []*ast.ImportSpec
	Vars    string
	Extends string

	Declarations string

	Handler string
}

func transformPageFile(projectPath string, filePath string) (result transformed) {
	by, err := os.ReadFile(filePath)
	if err != nil {
		panic(err)
	}
	result.RouterPath = pathIntoUrl(projectPath, filePath)
	result.Path = filePath

	content := string(by)
	commands := map[string]string{}

	htmlStart := 0
	end := findCodeBlock(string(by))
	if end != -1 {
		htmlStart = end + 3
		dirty := trim(content[3:end])

		code := ""

		for _, line := range strings.Split(dirty, "\n") {
			statement := trim(line)
			if len(statement) > 0 && statement[0] == '!' {
				// command statement
				statement = statement[1:]
				space := strings.Index(statement, " ")
				commands[statement[:space]] = statement[space+1:]
			} else {
				code += line + "\n"
			}
		}

		imports, importsEnd, err := getImports(code)
		if err != nil {
			panic(err)
		}
		result.Imports = imports
		result.RawCode = trim(code[importsEnd:])

		declarations, err := getDeclarations(result.RawCode)
		if err != nil {
			panic(err)
		}
		result.Declarations = makeDeclarations(declarations)

		result, err = applyCommands(commands, result)
		if err != nil {
			panic(err)
		}
	}

	result.Vars = fmt.Sprintf(
		"var __TEMPLATE = html.Must(html.New(\"__TEMPLATE\").Parse(\"%s\")) \n __TEMPLATE = html.Must(__TEMPLATE.Parse(\"%s\"))",
		escapeDoubleQuotes(removeNewlines(trim(result.Extends))),
		escapeDoubleQuotes(removeNewlines(trim(content[htmlStart:]))),
	)

	var tpl bytes.Buffer
	if err := handlerTemplate.Execute(&tpl, result); err != nil {
		panic(err)
	}
	result.Handler = tpl.String()

	return result
}

func getImports(code string) ([]*ast.ImportSpec, int, error) {
	offset := 0
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "", code, parser.PackageClauseOnly)
	if err != nil {
		prepend := "package handler\n"
		offset = -len(prepend)
		code = prepend + code
	}
	if file.Package.IsValid() {
		return nil, -1, errors.New("code block of a route handler may not contain a package declaration")
	}

	file, err = parser.ParseFile(fset, "", code, parser.ImportsOnly)
	if err != nil {
		return nil, -1, err
	}

	nChars := 0
	ast.Inspect(file, func(node ast.Node) bool {
		if genDecl, ok := node.(*ast.GenDecl); ok && genDecl.Tok == token.IMPORT {
			end := genDecl.End()
			endPos := fset.Position(end)
			lines := strings.Split(code, "\n")
			for i := 0; i <= endPos.Line-1 && i < len(lines); i++ {
				if i == endPos.Line-1 {
					nChars += endPos.Column
				} else {
					nChars += len(lines[i]) + 1
				}
			}
			return false
		}
		return true
	})

	if nChars != 0 {
		nChars = nChars + offset
	}

	return file.Imports, nChars, nil
}

func getDeclarations(handlerFunc string) ([]string, error) {
	handlerFunc = fmt.Sprintf("package yesverigut\nfunc Handler() {\n%s\n}", handlerFunc)

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "", handlerFunc, 0)
	if err != nil {
		return nil, err
	}

	declared := []string{}
	ast.Inspect(file, func(node ast.Node) bool {
		if fn, ok := node.(*ast.FuncDecl); ok {
			if fn.Name.Name != "Handler" {
				return true
			}

			for _, decl := range fn.Body.List {
				switch n := decl.(type) {
				case *ast.AssignStmt:
					for _, expr := range n.Lhs {
						if ident, ok := expr.(*ast.Ident); ok {
							if !slices.Contains(declared, ident.Name) {
								declared = append(declared, ident.Name)
							}
						}
					}
				case *ast.DeclStmt:
					if genDecl, ok := n.Decl.(*ast.GenDecl); ok {
						if genDecl.Tok == token.CONST || genDecl.Tok == token.VAR {
							for _, spec := range genDecl.Specs {
								if valueSpec, ok := spec.(*ast.ValueSpec); ok {
									for _, name := range valueSpec.Names {
										if !slices.Contains(declared, name.Name) {
											declared = append(declared, name.Name)
										}
									}
								}
							}
						}
					}
				}
			}
			return false
		}
		return true
	})

	return declared, nil
}

// transforms names of vars declared in the code block (array of strings)
// into a string of map entries, each in new line and separated with a comma
func makeDeclarations(declarationTags []string) (result string) {
	result = ""
	for _, tag := range declarationTags {
		result += fmt.Sprintf("\"%s\": %s,\n", tag, tag)
	}
	return result
}

// transforms slice of *ast.ImportSpec into a string of individual import statements
// search separated with a new line
func makeImports(imports []*ast.ImportSpec) (result string) {
	result = "import (\n"
	for _, spec := range imports {
		name := ""
		if spec.Name != nil {
			name = spec.Name.Name
		}
		result += fmt.Sprintf("%s %s\n", name, spec.Path.Value)
	}
	result += ")"
	return
}

// tries to find a code block in a component file, returns end index.
// If this index is equal to -1 then the component doesn't have a code block
func findCodeBlock(content string) int {
	if content[:3] != "---" {
		return -1
	}
	end := strings.Index(content[3:], "---")
	if end == -1 {
		return -1
	}

	return end + 3
}
