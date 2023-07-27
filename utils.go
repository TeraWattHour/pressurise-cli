package main

import (
	"flag"
	"strings"
)

func trim(str string) string {
	return strings.Trim(str, " 	\r\n")
}

func escapeDoubleQuotes(str string) string {
	return strings.ReplaceAll(str, "\"", "\\\"")
}

func escapeNewLines(str string) string {
	return strings.ReplaceAll(str, "\n", "\\n")
}

func filePathIntoRouterPath(path, appPath string) string {
	path = strings.TrimPrefix(path, appPath)
	path = strings.TrimSuffix(path, ".html")

	path = strings.TrimSuffix(path, "index")

	path = strings.TrimSuffix(path, "/")

	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	return path
}

func parseFlags(args []string) error {
	flagset := flag.CommandLine
	var positionalArgs []string
	for {
		if err := flagset.Parse(args); err != nil {
			return err
		}
		args = args[len(args)-flagset.NArg():]
		if len(args) == 0 {
			break
		}
		positionalArgs = append(positionalArgs, args[0])
		args = args[1:]
	}

	return flagset.Parse(positionalArgs)
}
