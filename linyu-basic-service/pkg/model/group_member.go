package model

import (
	"github.com/linyu-im/linyu-server/linyu-common/pkg/db"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/localtime"
	"gorm.io/gorm"
)

func init() {
	db.MysqlAddMigrateTable(&GroupMember{})
}

// GroupMember 聊天群成员表
type GroupMember struct {
	ID            string              `gorm:"size:64;primaryKey;autoIncrement:false;comment:id" json:"id"`
	GroupID       string              `gorm:"size:64;not null;comment:聊天群id" json:"groupId"`
	UserID        string              `gorm:"size:64;not null;comment:成员id" json:"userId"`
	GroupNickName string              `gorm:"size:128;comment:群昵称" json:"groupNickName"`
	GroupRemark   string              `gorm:"size:128;comment:群备注" json:"groupRemark"`
	CreatedAt     localtime.LocalTime `gorm:"type:timestamp(3);not null;autoCreateTime;comment:创建时间" json:"createdAt"`
	UpdatedAt     localtime.LocalTime `gorm:"type:timestamp(3);not null;autoUpdateTime;comment:更新时间" json:"updatedAt"`
	DeletedAt     gorm.DeletedAt      `gorm:"index" json:"deletedAt"`
}

func (GroupMember) TableName() string {
	return "t_group_member"
}

func (GroupMember) TableComment() string {
	return "聊天群成员表"
}
