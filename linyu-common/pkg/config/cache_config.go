package config

type CacheType string

const (
	RedisCacheType CacheType = "redis"
	OtterCacheType CacheType = "otter" // 本地缓存
)

type RedisCacheConfig struct {
	Addr     string `mapstructure:"addr"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

type OtterConfig struct {
	Capacity int `mapstructure:"capacity"`
}
