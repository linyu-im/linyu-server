package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Server struct {
		Port        int    `mapstructure:"port"`
		RoutePrefix string `mapstructure:"route-prefix"`
		NodeId      int64  `mapstructure:"node-id"`
	} `mapstructure:"server"`

	Rdb struct {
		Type      RdbType         `mapstructure:"type"`
		MysqlRdb  MysqlRdbConfig  `mapstructure:"mysql"`
		SqliteRdb SqliteRdbConfig `mapstructure:"sqlite"`
	} `mapstructure:"rdb"`

	EventBus struct {
		Type           EventBusType         `mapstructure:"type"`
		RocketEventBus RocketEventBusConfig `mapstructure:"rocket"`
	} `mapstructure:"event-bus"`

	Cache struct {
		Type       CacheType        `mapstructure:"type"`
		RedisCache RedisCacheConfig `mapstructure:"redis"`
		OtterCache OtterConfig      `mapstructure:"otter"`
	} `mapstructure:"cache"`

	Vector struct {
		Type           VectorType           `mapstructure:"type"`
		ChromemVector  ChromemVectorConfig  `mapstructure:"chromem"`
		WeaviateVector WeaviateVectorConfig `mapstructure:"weaviate"`
	} `mapstructure:"vector"`

	Jwt struct {
		Secret      string `mapstructure:"secret"`
		ExpireHours int    `mapstructure:"expire-hours"`
	} `mapstructure:"jwt"`

	Email struct {
		Host     string `mapstructure:"host"`
		Port     int    `mapstructure:"port"`
		FormAddr string `mapstructure:"form-addr"`
		Password string `mapstructure:"password"`
		Username string `mapstructure:"username"`
	} `mapstructure:"email"`

	Storage struct {
		Type         StorageType        `mapstructure:"type"`
		LocalStorage LocalStorageConfig `mapstructure:"local"`
		OssStorage   OssStorageConfig   `mapstructure:"oss"`
		S3Storage    S3StorageConfig    `mapstructure:"s3"`
	} `mapstructure:"storage"`

	Livekit struct {
		Host   string `mapstructure:"host"`
		Key    string `mapstructure:"key"`
		Secret string `mapstructure:"secret"`
	} `mapstructure:"livekit"`

	Auth struct {
		Ldap  LdapConfig  `mapstructure:"ldap"`
		Gitee GiteeConfig `mapstructure:"gitee"`
	} `mapstructure:"auth"`

	AI struct {
		Embedding EmbeddingConfig `mapstructure:"embedding"`
	} `mapstructure:"ai"`
}

var C Config

func InitConfig(path string) {
	v := viper.New()
	v.SetConfigFile(path)

	if err := v.ReadInConfig(); err != nil {
		panic("failed to read config file:" + err.Error())
	}

	if err := v.Unmarshal(&C); err != nil {
		panic("failed to read config file:" + err.Error())
	}
}
