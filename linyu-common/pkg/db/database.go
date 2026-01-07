package db

import (
	"github.com/linyu-im/linyu-server/linyu-common/pkg/db/mysql"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/db/redis"
	"gorm.io/gorm"
)

var MysqlDB *gorm.DB
var RedisDB *redis.RedisClient

func Init() {
	MysqlDB = mysql.CreateMysqlDB()
	RedisDB = redis.CreateRedisClient()
	//表结构新建
	err := MysqlMigrate(MysqlDB, mysqlModels)
	if err != nil {
		panic("mysql db migrate error: " + err.Error())
	}
}

var mysqlModels []interface{}

func MysqlAddMigrateTable(model interface{}) {
	mysqlModels = append(mysqlModels, model)
}

func MysqlMigrate(MysqlDB *gorm.DB, models []interface{}) error {
	for _, m := range models {
		comment := ""
		if tc, ok := m.(interface{ TableComment() string }); ok {
			comment = tc.TableComment()
		}
		if err := MysqlDB.Set("gorm:table_options", "COMMENT='"+comment+"'").AutoMigrate(m); err != nil {
			return err
		}
	}
	return nil
}
