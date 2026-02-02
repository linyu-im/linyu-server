package rdb

import (
	"github.com/linyu-im/linyu-server/linyu-common/pkg/config"
	"gorm.io/gorm"
)

func InitRdb(models []interface{}) *gorm.DB {
	if config.C.Rdb.Type == "" {
		panic("rdb type not set")
	}
	var rdb *gorm.DB
	switch config.C.Rdb.Type {
	case config.MysqlRdbType:
		rdb = NewCreateMysqlDB()
		err := MysqlMigrate(rdb, models)
		if err != nil {
			panic("mysql db migrate error: " + err.Error())
		}
	case config.SqliteRdbType:
		rdb = NewCreateSqliteDB()
		err := SqliteMigrate(rdb, models)
		if err != nil {
			panic("sqlite db migrate error: " + err.Error())
		}
	default:
		panic("storage type not supported: " + config.C.Storage.Type)
	}
	return rdb
}
