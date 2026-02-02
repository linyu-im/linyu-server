package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Server struct {
		Port        int    `mapstructure:"port"`
		RoutePrefix string `mapstructure:"route-prefix"`
	} `mapstructure:"server"`

	Rdb struct {
		Type      RdbType         `mapstructure:"type"`
		MysqlRdb  MysqlRdbConfig  `mapstructure:"mysql"`
		SqliteRdb SqliteRdbConfig `mapstructure:"sqlite"`
	} `mapstructure:"rdb"`

	Cache struct {
		Type       CacheType        `mapstructure:"type"`
		RedisCache RedisCacheConfig `mapstructure:"redis"`
		OtterCache OtterConfig      `mapstructure:"otter"`
	} `mapstructure:"cache"`

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
	}
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
