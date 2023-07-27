package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func printUsage() {
	prompt := "Wrong arguments provided.\n\n" +
		"Usage:\n" +
		"  pressurise-cli build <project-directory> [<optional-arguments>]"

	fmt.Println(prompt)

	if flag.Parsed() {
		fmt.Println("\nOptional arguments:")

		flag.PrintDefaults()
	}
}

func main() {
	help := flag.Bool("help", false, "a `bool`, whether to show help")
	args := os.Args[1:]

	if len(args) == 0 {
		printUsage()
		return
	}

	mode := args[0]
	if mode == "build" && len(args) >= 2 {
		noformat := flag.Bool("no-format", false, "a `bool`, whether output file should be formatted with go fmt")

		parseFlags(args[2:])

		if *help {
			printUsage()
			return
		}

		projectDirectory := args[1]

		if err := buildProject(projectDirectory); err != nil {
			fmt.Println("project couldn't be built:", err)
			return
		}

		if !*noformat {
			formatGenerated := exec.Command("go", "fmt", filepath.Join(projectDirectory, "pressrelease.go"))
			formatGenerated.Start()
		}
	} else {
		printUsage()
	}
}
