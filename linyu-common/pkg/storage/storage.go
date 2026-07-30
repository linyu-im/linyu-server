package storage

import (
	"github.com/linyu-im/linyu-server/linyu-common/pkg/config"
	"io"
)

var S Storage

// mergeUploadConcurrency 合并时并发上传分片数，避免串行二次上传过慢
const mergeUploadConcurrency = 8

type Storage interface {
	Upload(fileKey string, reader io.Reader) (string, error)
	Download(fileKey string, writer io.Writer) error
	Delete(fileKey string) error
	Merge(fileKey string, chunkDir string, totalChunks int) (string, error)
}

func InitStorage() {
	if config.C.Storage.Type == "" {
		panic("storage type not set")
	}
	switch config.C.Storage.Type {
	case config.LocalStorageType:
		S = NewLocalStorage(config.C.Storage.LocalStorage)
	case config.OssStorageType:
		S = NewOssStorage(config.C.Storage.OssStorage)
	case config.S3StorageType:
		S = NewS3Storage(config.C.Storage.S3Storage)
	default:
		panic("storage type not supported: " + config.C.Storage.Type)
	}
}
