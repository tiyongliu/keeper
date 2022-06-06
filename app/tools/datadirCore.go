package tools

import (
	"os"
	"strings"
)

const MkDir = "keeper-data"

func DataDirCore() string {
	dir, _ := os.UserHomeDir()
	return strings.Join([]string{dir, MkDir}, "/")
}
