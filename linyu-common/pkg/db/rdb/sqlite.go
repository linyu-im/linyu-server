package rdb

import (
	"github.com/glebarez/sqlite"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/config"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"os"
	"path/filepath"
)

func NewCreateSqliteDB() *gorm.DB {
	c := config.C.Rdb.SqliteRdb
	if _, err := os.Stat(c.DbPath); os.IsNotExist(err) {
		_ = os.MkdirAll(c.DbPath, os.ModePerm)
	}
	dbFile := filepath.Join(c.DbPath, c.DbName+".db")
	db, err := gorm.Open(sqlite.Open(dbFile), &gorm.Config{})
	if err != nil {
		panic("failed to connect sqlite database: " + err.Error())
	}

	db.Logger.LogMode(logger.Warn)
	return db
}

func SqliteMigrate(SqliteDB *gorm.DB, models []interface{}) error {
	for _, m := range models {
		if err := SqliteDB.AutoMigrate(m); err != nil {
			return err
		}
	}
	return nil
}
