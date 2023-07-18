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

func pathIntoUrl(projectPath string, filePath string) string {
	u := strings.TrimSuffix(strings.TrimSuffix(strings.TrimPrefix(filePath, projectPath), ".html"), "index")
	if len(u) == 0 {
		return "/"
	}
	if u[0] != '/' {
		u = "/" + u
	}
	if len(u) > 1 {
		u = strings.TrimSuffix(u, "/")
	}
	return u
}
