package storage

import (
	"github.com/linyu-im/linyu-server/linyu-common/pkg/config"
	"io"
)

var S Storage

type Storage interface {
	Upload(fileKey string, reader io.Reader) (string, error)
	Download(fileKey string, writer io.Writer) error
	Delete(fileKey string) error
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
	default:
		panic("storage type not supported: " + config.C.Storage.Type)
	}
}
