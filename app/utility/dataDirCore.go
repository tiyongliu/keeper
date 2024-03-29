package utility

import (
	"keeper/app/pkg/logger"
	"os"
	"path"
)

const MkDir = ".keeper"

func ensureDirectory(dir string) {
	if !IsExist(dir) {
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			logger.Fatalf("os.MkdirAll failed err: %v\n", err)
		}
	}
}
func DataDirCore() string {
	dir, _ := os.UserHomeDir()
	return path.Join(dir, MkDir, "/")
}

func DataDir() string {
	dir := DataDirCore()
	ensureDirectory(dir)
	return dir
}
