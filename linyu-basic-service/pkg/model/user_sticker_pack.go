package model

import (
	"github.com/linyu-im/linyu-server/linyu-common/pkg/db"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/localtime"
	"gorm.io/gorm"
)

func init() {
	db.AddMigrateTable(&UserStickerPack{})
}

// UserStickerPack 用户拥有的表情分组表
type UserStickerPack struct {
	ID            string              `gorm:"size:64;primaryKey;autoIncrement:false;comment:id" json:"id"`
	UserID        string              `gorm:"size:64;index;not null;comment:用户id" json:"userId"`
	StickerPackID string              `gorm:"size:64;index;not null;comment:表情分组id" json:"stickerPackId"`
	CreatedAt     localtime.LocalTime `gorm:"type:timestamp(3);not null;autoCreateTime;comment:创建时间" json:"createdAt"`
	UpdatedAt     localtime.LocalTime `gorm:"type:timestamp(3);not null;autoUpdateTime;comment:更新时间" json:"updatedAt"`
	DeletedAt     gorm.DeletedAt      `gorm:"index" json:"deletedAt"`
}

func (UserStickerPack) TableName() string {
	return "t_user_sticker_pack"
}

func (UserStickerPack) TableComment() string {
	return "用户拥有的表情分组表"
}
