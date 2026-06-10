package model

import (
	"github.com/linyu-im/linyu-server/linyu-common/pkg/db"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/localtime"
	"gorm.io/gorm"
)

func init() {
	db.AddMigrateTable(&UserStickerFavorite{})
}

// UserStickerFavorite 用户收藏的表情表
type UserStickerFavorite struct {
	ID        string              `gorm:"size:64;primaryKey;autoIncrement:false;comment:id" json:"id"`
	UserID    string              `gorm:"size:64;index;not null;comment:用户id" json:"userId"`
	StickerID string              `gorm:"size:64;index;not null;comment:表情id" json:"stickerId"`
	CreatedAt localtime.LocalTime `gorm:"type:timestamp(3);not null;autoCreateTime;comment:创建时间" json:"createdAt"`
	UpdatedAt localtime.LocalTime `gorm:"type:timestamp(3);not null;autoUpdateTime;comment:更新时间" json:"updatedAt"`
	DeletedAt gorm.DeletedAt      `gorm:"index" json:"deletedAt"`
}

func (UserStickerFavorite) TableName() string {
	return "t_user_sticker_favorite"
}

func (UserStickerFavorite) TableComment() string {
	return "用户收藏的表情表"
}
