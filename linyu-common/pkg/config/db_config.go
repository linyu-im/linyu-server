package config

//---------------------关系型---------------------

type RdbType string

const (
	MysqlRdbType  RdbType = "mysql"
	SqliteRdbType RdbType = "sqlite"
)

type MysqlRdbConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Database string `mapstructure:"database"`
	Timezone string `mapstructure:"timezone"`
}

type SqliteRdbConfig struct {
	DbPath string `mapstructure:"db-path"`
	DbName string `mapstructure:"db-name"`
}

//---------------------缓存---------------------

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

//---------------------向量---------------------

type VectorType string

const (
	ChromemVectorType  VectorType = "chromem" //本地向量
	WeaviateVectorType VectorType = "weaviate"
)

type ChromemVectorConfig struct {
	Path string `mapstructure:"path"`
}

type WeaviateVectorConfig struct {
	Scheme string `mapstructure:"scheme"`
	Host   string `mapstructure:"host"`
}
