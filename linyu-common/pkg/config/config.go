package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Server struct {
		Port int `mapstructure:"port"`
	} `mapstructure:"server"`

	Mysql struct {
		Host     string `mapstructure:"host"`
		Port     int    `mapstructure:"port"`
		User     string `mapstructure:"user"`
		Password string `mapstructure:"password"`
		Database string `mapstructure:"database"`
		Timezone string `mapstructure:"timezone"`
	} `mapstructure:"mysql"`

	Redis struct {
		Addr     string `mapstructure:"addr"`
		Password string `mapstructure:"password"`
		DB       int    `mapstructure:"db"`
	} `mapstructure:"redis"`

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
