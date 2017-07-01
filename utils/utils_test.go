package utils

import (
	"fmt"
	"testing"
)

func TestGetGOPATHs(t *testing.T) {
	path := GetGoPATHs()

	if len(path) == 0 {
		t.Error("GOPATH environment variable is not set or empty")
	}

	fmt.Println("path--> ", path)
}
