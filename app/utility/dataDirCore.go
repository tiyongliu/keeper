package utility

import (
	"log"
	"os"
	"strings"
)

const MkDir = ".keeper"

func ensureDirectory(dir string, clean bool) {
	if !IsExist(dir) {
		if err := os.MkdirAll(dir, SecondFilePerm); err != nil {
			log.Fatalf("os.MkdirAll failed err: %v\n", err)
		}
	}
}
func DataDirCore() string {
	dir, _ := os.UserHomeDir()
	return strings.Join([]string{dir, MkDir}, "/")
}

func DataDir() string {
	dir := DataDirCore()
	ensureDirectory(dir, false)
	return dir
}
