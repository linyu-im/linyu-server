package model

import (
	"github.com/linyu-im/linyu-server/linyu-common/pkg/db"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/localtime"
	"gorm.io/gorm"
)

func init() {
	db.MysqlAddMigrateTable(&PhysicalFile{})
}

type PhysicalFile struct {
	ID          string              `gorm:"size:64;primaryKey;autoIncrement:false;comment:id" json:"id"`
	FileHash    string              `gorm:"size:64;uniqueIndex;not null;comment:文件hash" json:"fileHash"`
	FileSize    int64               `gorm:"type:bigint;not null;comment:文件大小" json:"fileSize"`
	StoragePath string              `gorm:"type:text;not null;comment:存储路径" json:"storagePath"`
	RefCount    int                 `gorm:"type:int;default:1;comment:引用计数" json:"refCount"`
	CreatedAt   localtime.LocalTime `gorm:"type:timestamp(3);autoCreateTime"`
	UpdatedAt   localtime.LocalTime `gorm:"type:timestamp(3);autoUpdateTime"`
	DeletedAt   gorm.DeletedAt      `gorm:"index"`
}

func (PhysicalFile) TableName() string {
	return "t_physical_file"
}

func (PhysicalFile) TableComment() string {
	return "物理文件表"
}
