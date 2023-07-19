package main

import (
	"strings"
)

func trim(str string) string {
	return strings.Trim(str, " 	\r\n")
}

func escapeDoubleQuotes(str string) string {
	return strings.ReplaceAll(str, "\"", "\\\"")
}

func removeNewlines(str string) string {
	return strings.ReplaceAll(str, "\n", "")
}

func pathIntoUrl(projectPath string, filePath string) (x string) {
	x = strings.TrimPrefix(filePath, projectPath)
	x = strings.TrimSuffix(x, ".html")
	x = strings.TrimSuffix(x, "index")
	if len(x) == 0 {
		return "/"
	}
	if x[0] != '/' {
		x = "/" + x
	}
	if len(x) > 1 {
		x = strings.TrimSuffix(x, "/")
	}
	return x
}
