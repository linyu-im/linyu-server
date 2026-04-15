package rdb

import (
	"fmt"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/config"
	logger2 "github.com/linyu-im/linyu-server/linyu-common/pkg/logger"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewCreateMysqlDB() *gorm.DB {
	c := config.C.Rdb.MysqlRdb
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.User, c.Password, c.Host, c.Port, c.Database)
	var err error
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	db.Logger.LogMode(logger.Warn)
	if err != nil {
		panic("failed to connect mysql database: " + err.Error())
	}
	return db
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
		if tc, ok := m.(interface{ TableInit(db *gorm.DB) error }); ok {
			err := tc.TableInit(MysqlDB)
			if err != nil {
				logger2.Log.Error("[MysqlMigrate] 数据初始化失败:", zap.Error(err))
			}
		}
	}
	return nil
}
