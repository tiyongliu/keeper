package tools

import (
	"os"
	"strings"
)

const MkDir = ".keeper"

func DataDirCore() string {
	dir, _ := os.UserHomeDir()
	return strings.Join([]string{dir, MkDir}, "/")
}
