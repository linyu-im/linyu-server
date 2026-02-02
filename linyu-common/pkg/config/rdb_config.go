package config

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
