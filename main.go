package main

import (
	"fmt"
	"os"
)

func printUsage() {
	prompt := "Wrong arguments provided.\n\n" +
		"Usage:\n" +
		"	pressurise-cli build <project-directory>"

	fmt.Println(prompt)
}

func main() {
	args := os.Args[1:]

	mode := args[0]

	if mode == "build" && len(args) >= 2 {
		projectDirectory := args[1]

		if err := buildProject(projectDirectory); err != nil {
			fmt.Println("project couldn't be built:", err)
			return
		}

	} else {
		printUsage()
	}
}
