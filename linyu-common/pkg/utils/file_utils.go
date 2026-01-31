package utils

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

func ResolvePath(filePath, fileKey string) (string, error) {
	absBasePath, err := filepath.Abs(filePath)
	if err != nil {
		return "", err
	}

	fullPath := filepath.Join(absBasePath, fileKey)
	rel, err := filepath.Rel(absBasePath, fullPath)
	if err != nil {
		return "", errors.New("unsafe file path detected")
	}

	if strings.HasPrefix(rel, "..") || rel == ".." {
		return "", errors.New("unsafe file path detected")
	}

	return fullPath, nil
}

func GetFileSystemPath(filePath, fileKey string) (string, error) {
	fullPath, err := ResolvePath(filePath, fileKey)
	if err != nil {
		return "", err
	}
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return "", errors.New("file does not exist")
	}
	return fullPath, nil
}
