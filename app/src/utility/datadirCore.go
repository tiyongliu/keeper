package utility

import (
	"os"
	"strings"
)

const MkDir = "vaults-data"

func dataDirCore() string {
	dir, _ := os.UserHomeDir()
	return strings.Join([]string{dir, MkDir}, "/")
}
