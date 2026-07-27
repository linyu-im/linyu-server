package model

import (
	"github.com/linyu-im/linyu-server/linyu-common/pkg/db"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/localtime"
)

func init() {
	db.AddMigrateTable(&UserBadge{})
}

// UserBadge 用户红点表
type UserBadge struct {
	UserID     string              `gorm:"size:64;primaryKey;comment:用户id" json:"userId"`
	BadgeCode  string              `gorm:"size:100;primaryKey;comment:红点编码" json:"badgeCode"`
	LastReadID string              `gorm:"size:64;not null;default:'';comment:最后已读id" json:"lastReadId"`
	UpdatedAt  localtime.LocalTime `gorm:"type:timestamp(3);not null;autoUpdateTime;comment:更新时间" json:"updatedAt"`

	UnreadCount int `gorm:"->;-:migration" json:"unreadCount"`
}

func (UserBadge) TableName() string {
	return "t_user_badge"
}

func (UserBadge) TableComment() string {
	return "用户红点表"
}
