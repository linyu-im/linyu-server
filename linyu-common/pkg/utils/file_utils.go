package utils

import (
	"errors"
	"io"
	"io/fs"
	"mime/multipart"
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

func SaveUploadedFile(file *multipart.FileHeader, dst string, perm ...fs.FileMode) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	var mode os.FileMode = 0o750
	if len(perm) > 0 {
		mode = perm[0]
	}
	dir := filepath.Dir(dst)
	if err = os.MkdirAll(dir, mode); err != nil {
		return err
	}
	if err = os.Chmod(dir, mode); err != nil {
		return err
	}

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	return err
}
