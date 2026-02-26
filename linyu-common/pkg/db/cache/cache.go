package cache

import (
	"github.com/linyu-im/linyu-server/linyu-common/pkg/config"
	"time"
)

type Cache interface {
	Set(key string, value interface{}, expiration time.Duration) error
	Get(key string) (string, error)
	SetObject(key string, value interface{}, expiration time.Duration) error
	GetObject(key string, dest interface{}) error
	Del(key string) error
	Exists(key string) (bool, error)
	SAdd(key string, expiration time.Duration, members ...interface{}) error
	SMembers(key string) ([]string, error)
}

func InitCache() Cache {
	if config.C.Cache.Type == "" {
		panic("cache type not set")
	}
	var cache Cache
	switch config.C.Cache.Type {
	case config.RedisCacheType:
		cache = NewRedisClient()
	case config.OtterCacheType:
		cache = NewOtterClient()
	default:
		panic("storage type not supported: " + config.C.Storage.Type)
	}
	return cache
}
