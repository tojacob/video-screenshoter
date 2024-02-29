package utils

import (
	"fmt"
	"os"
	"path/filepath"
)

func CheckError(err error) {
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
}

func GetCurrentPath() string {
	dir, err := os.Getwd()
	CheckError(err)

	return dir
}

func GetRelativePath(basePath string, fullPath string) string {
	relativePath, err := filepath.Rel(basePath, fullPath)
	CheckError(err)

	return relativePath
}
