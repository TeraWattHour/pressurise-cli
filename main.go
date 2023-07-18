package main

import (
	"fmt"
	"os"
	"os/exec"
	"path"
)

func printUsage() {
	prompt := "Wrong arguments provided.\n\n" +
		"Usage:\n" +
		"	pressurise-cli build <project-directory>"

	fmt.Println(prompt)
}

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		printUsage()
		return
	}

	mode := args[0]

	if mode == "build" && len(args) >= 2 {
		projectDirectory := args[1]
		formatGenerated := exec.Command("go", "fmt", path.Join(projectDirectory, "pressrelease.go"))

		if err := buildProject(projectDirectory); err != nil {
			fmt.Println("project couldn't be built:", err)
			return
		}
		formatGenerated.Start()
	} else {
		printUsage()
	}
}
