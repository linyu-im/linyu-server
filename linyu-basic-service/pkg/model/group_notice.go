package model

import (
	"github.com/linyu-im/linyu-server/linyu-common/pkg/db"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/localtime"
	"gorm.io/gorm"
)

func init() {
	db.AddMigrateTable(&GroupNotice{})
}

// GroupNotice 群公告
type GroupNotice struct {
	ID              string              `gorm:"size:64;primaryKey;autoIncrement:false;comment:群公告id" json:"id"`
	GroupID         string              `gorm:"size:64;not null;index:idx_group_notice_group_id;comment:群id" json:"groupId"`
	PublisherUserID string              `gorm:"size:64;not null;index:idx_group_notice_publisher_user_id;comment:发布人用户id" json:"publisherUserId"`
	Content         string              `gorm:"type:text;not null;comment:公告内容" json:"content"`
	IsTop           bool                `gorm:"not null;default:false;index:idx_group_notice_is_top;comment:是否置顶" json:"isTop"`
	CreatedAt       localtime.LocalTime `gorm:"type:timestamp(3);not null;autoCreateTime;comment:创建时间" json:"createdAt"`
	UpdatedAt       localtime.LocalTime `gorm:"type:timestamp(3);not null;autoUpdateTime;comment:更新时间" json:"updatedAt"`
	DeletedAt       gorm.DeletedAt      `gorm:"index" json:"deletedAt"`
}

func (GroupNotice) TableName() string {
	return "t_group_notice"
}

func (GroupNotice) TableComment() string {
	return "群公告表"
}
