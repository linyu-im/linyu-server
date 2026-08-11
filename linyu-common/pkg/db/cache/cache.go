package cache

import (
	"time"

	"github.com/linyu-im/linyu-server/linyu-common/pkg/config"
)

type Cache interface {
	Set(key string, value interface{}, expiration time.Duration) error
	Get(key string) (string, error)
	SetObject(key string, value interface{}, expiration time.Duration) error
	GetObject(key string, dest interface{}) error
	Del(key string) error
	Exists(key string) (bool, error)
	SAdd(key string, expiration time.Duration, members ...interface{}) error
	SRem(key string, members ...interface{}) error
	SMembers(key string) ([]string, error)
	// SetIfGreater 仅当 newValue 数值大于当前值（或不存在）时写入，用于分布式下安全更新最大值
	SetIfGreater(key string, newValue string, expiration time.Duration) error
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
