package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
)

func printUsage() {
	prompt :=
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

	if len(args) < 2 {
		printUsage()
		return
	}

	mode := args[0]
	args = args[1:]
	if mode == "build" {
		noformat := flag.Bool("no-format", false, "a `bool`, whether output file should be formatted with go fmt")

		parseFlags(args[1:])

		if *help {
			printUsage()
			return
		}

		projectDirectory := args[0]

		outputFile, err := buildProject(projectDirectory)
		if err != nil {
			fmt.Printf("\033[31m██\033[0m Project couldn't be built:\n%s\n", err)
			return
		}
		fmt.Println("\033[32m██\033[0m Project built successfully")

		if !*noformat {
			exec.Command("go", "fmt", outputFile).Start()
		} else {
			fmt.Println("\033[38;5;208m██\033[0m Skipping formatting")
		}
	} else {
		printUsage()
	}
}
