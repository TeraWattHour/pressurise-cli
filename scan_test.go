package main

import (
	"testing"

	"golang.org/x/exp/slices"
)

// very nice and comprehensive test xdd
func TestScanPagesDirectory(t *testing.T) {
	found, err := scanPagesDirectory("tests/basic/app/", []string{})
	if err != nil {
		t.Error(err)
	}

	if len(found) != 2 || !slices.Contains(found, "tests/basic/app/blob/index.html") || !slices.Contains(found, "tests/basic/app/index.html") {
		t.Error("bad structure")
	}
}

func TestVerifyProjectDirectory(t *testing.T) {
	err := verifyProjectDirectory("tests/basic/")
	if err != nil {
		t.Error(err)
	}

	err = verifyProjectDirectory("tests/")
	if err == nil {
		t.Error("expected directory not to be a valid project", err)
	}
}
