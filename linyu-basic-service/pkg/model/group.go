package model

import (
	"github.com/linyu-im/linyu-server/linyu-common/pkg/db"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/localtime"
	"gorm.io/gorm"
)

func init() {
	db.MysqlAddMigrateTable(&Group{})
}

// Group 聊天群
type Group struct {
	ID            string              `gorm:"size:64;primaryKey;autoIncrement:false;comment:群id" json:"id"`
	CreatorUserID string              `gorm:"size:64;not null;comment:创建用户id" json:"creatorUserId"`
	GroupNumber   string              `gorm:"size:64;not null;comment:群号" json:"groupNumber"`
	OwnerUserID   string              `gorm:"size:64;not null;comment:群主用户id" json:"ownerUserId"`
	Name          string              `gorm:"size:128;not null;comment:群名名称" json:"name"`
	Avatar        string              `gorm:"size:512;comment:群头像URL" json:"avatar"`
	Describe      string              `gorm:"type:text;comment:群描述" json:"describe"`
	MemberNum     int                 `gorm:"default:0;comment:成员数" json:"memberNum"`
	CreatedAt     localtime.LocalTime `gorm:"type:timestamp(3);not null;autoCreateTime;comment:创建时间" json:"createdAt"`
	UpdatedAt     localtime.LocalTime `gorm:"type:timestamp(3);not null;autoUpdateTime;comment:更新时间" json:"updatedAt"`
	DeletedAt     gorm.DeletedAt      `gorm:"index" json:"deletedAt"`
}

func (Group) TableName() string {
	return "t_group"
}

func (Group) TableComment() string {
	return "聊天群表"
}
