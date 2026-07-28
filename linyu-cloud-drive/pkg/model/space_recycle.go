package model

import (
	"github.com/linyu-im/linyu-server/linyu-common/pkg/db"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/localtime"
	"gorm.io/gorm"
)

func init() {
	db.AddMigrateTable(&SpaceRecycle{})
}

// SpaceRecycle 空间回收站表
type SpaceRecycle struct {
	ID          string               `gorm:"size:64;primaryKey;autoIncrement:false;comment:id" json:"id"`
	UserID      string               `gorm:"size:64;index;not null;comment:所属用户id" json:"userId"`
	SpaceID     string               `gorm:"size:64;index;not null;comment:空间id" json:"spaceId"`
	SpaceFileID string               `gorm:"size:64;uniqueIndex;not null;comment:空间文件id" json:"spaceFileId"`
	DeletedBy   string               `gorm:"size:64;index;not null;comment:删除人用户id" json:"deletedBy"`
	ExpireAt    *localtime.LocalTime `gorm:"type:timestamp(3);index;comment:过期时间，到期可彻底删除" json:"expireAt"`
	CreatedAt   localtime.LocalTime  `gorm:"type:timestamp(3);autoCreateTime;comment:创建时间" json:"createdAt"`
	UpdatedAt   localtime.LocalTime  `gorm:"type:timestamp(3);autoUpdateTime;comment:更新时间" json:"updatedAt"`
	DeletedAt   gorm.DeletedAt       `gorm:"index" json:"deletedAt"`

	// 关联 SpaceFile 字段
	FileName string `gorm:"->;-:migration" json:"fileName"`
	IsDir    bool   `gorm:"->;-:migration" json:"isDir"`
	FileType string `gorm:"->;-:migration" json:"fileType"`
	FileSize int64  `gorm:"->;-:migration" json:"fileSize"`
}

func (SpaceRecycle) TableName() string {
	return "t_space_recycle"
}

func (SpaceRecycle) TableComment() string {
	return "空间回收站表"
}
