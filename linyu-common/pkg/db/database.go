package db

import (
	"github.com/linyu-im/linyu-server/linyu-common/pkg/db/cache"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/db/rdb"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/db/vector"
	"gorm.io/gorm"
)

var RDB *gorm.DB         //关系型
var CacheDB cache.Cache  //缓存
var Vector vector.Vector //向量

func InitDB() {
	RDB = rdb.InitRdb(models)
	CacheDB = cache.InitCache()
	Vector = vector.InitVector()
}

var models []interface{}

func AddMigrateTable(model interface{}) {
	models = append(models, model)
}
