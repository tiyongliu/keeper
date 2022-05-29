package utility

import (
	"os"
	"strings"
)

const MkDir = "keeper-data"

func dataDirCore() string {
	dir, _ := os.UserHomeDir()
	return strings.Join([]string{dir, MkDir}, "/")
}
