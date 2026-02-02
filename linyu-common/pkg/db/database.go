package db

import (
	"github.com/linyu-im/linyu-server/linyu-common/pkg/db/rdb"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/db/redis"
	"gorm.io/gorm"
)

var RDB *gorm.DB
var RedisDB *redis.RedisClient

func InitDB() {
	RDB = rdb.InitRdb(mysqlModels)
	RedisDB = redis.CreateRedisClient()
}

var mysqlModels []interface{}

func MysqlAddMigrateTable(model interface{}) {
	mysqlModels = append(mysqlModels, model)
}
