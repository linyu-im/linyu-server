package storage

import (
	"errors"
	"fmt"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/config"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/utils"
	"io"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
)

var LocalStorageUrl = "/basic/v1/local/storage"

func NewLocalStorage(c config.LocalStorageConfig) *LocalStorage {
	if c.BaseURL == "" {
		panic("local storage required config: BaseURL not set")
	}
	filePath := c.FilePath
	if filePath == "" {
		wd, err := os.Getwd()
		if err != nil {
			panic("failed to get current working directory for LocalStorage FilePath: " + err.Error())
		}
		filePath = filepath.Join(wd, "linyu-storage")
	}
	if err := os.MkdirAll(filePath, 0755); err != nil {
		panic("failed to create storage directory: " + err.Error())
	}
	return &LocalStorage{
		BaseURL:  c.BaseURL,
		FilePath: c.FilePath,
	}
}

type LocalStorage struct {
	BaseURL  string
	FilePath string
}

func (s *LocalStorage) Upload(fileKey string, reader io.Reader) (string, error) {
	fullPath, err := utils.ResolvePath(s.FilePath, fileKey)
	if err != nil {
		return "", err
	}

	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", errors.New("failed to create storage directory")
	}

	out, err := os.OpenFile(fullPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return "", errors.New("failed to create target file")
	}
	defer out.Close()

	_, err = io.Copy(out, reader)
	if err != nil {
		return "", errors.New("failed to save file content")
	}

	return s.GetURL(fileKey, 0)
}

func (s *LocalStorage) Download(fileKey string, writer io.Writer) error {
	fullPath, err := utils.ResolvePath(s.FilePath, fileKey)
	if err != nil {
		return err
	}
	file, err := os.Open(fullPath)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = io.Copy(writer, file)
	return err
}

func (s *LocalStorage) Delete(fileKey string) error {
	if fileKey == "" {
		return nil
	}

	fullPath, err := utils.ResolvePath(s.FilePath, fileKey)
	if err != nil {
		return err
	}

	fi, err := os.Stat(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	if fi.IsDir() {
		return errors.New("cannot delete a directory via file storage interface")
	}

	err = os.Remove(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	return nil
}

func (s *LocalStorage) GetURL(fileKey string, expire int64) (string, error) {
	if fileKey == "" {
		return "", nil
	}

	u, err := url.Parse(s.BaseURL + config.C.Server.RoutePrefix + LocalStorageUrl)
	if err != nil {
		return "", err
	}

	u.Path = path.Join(u.Path, fileKey)

	if strings.HasSuffix(fileKey, "/") && !strings.HasSuffix(u.Path, "/") {
		u.Path += "/"
	}

	return u.String(), nil
}

func (s *LocalStorage) Merge(fileKey string, chunkDir string, totalChunks int) (string, error) {
	fullPath, err := utils.ResolvePath(s.FilePath, fileKey)
	if err != nil {
		return "", err
	}

	if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
		return "", err
	}

	targetFile, err := os.OpenFile(fullPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return "", err
	}
	defer targetFile.Close()

	for i := 0; i < totalChunks; i++ {
		partPath := fmt.Sprintf("%s/%d.part", chunkDir, i)
		partData, err := os.Open(partPath)
		if err != nil {
			return "", err
		}
		_, err = io.Copy(targetFile, partData)
		partData.Close()
		if err != nil {
			return "", err
		}
	}

	return s.GetURL(fileKey, 0)
}
