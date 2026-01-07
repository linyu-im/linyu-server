package mysql

import (
	"fmt"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// CreateMysqlDB 创建Mysql客户端
func CreateMysqlDB() *gorm.DB {
	c := config.C.Mysql
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.User, c.Password, c.Host, c.Port, c.Database)
	var err error
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect mysql database: " + err.Error())
	}
	return db
}
