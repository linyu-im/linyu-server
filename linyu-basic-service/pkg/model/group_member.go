package model

import (
	"github.com/linyu-im/linyu-server/linyu-common/pkg/db"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/localtime"
	"gorm.io/gorm"
)

func init() {
	db.AddMigrateTable(&GroupMember{})
}

// GroupMember 聊天群成员表
type GroupMember struct {
	ID             string              `gorm:"size:64;primaryKey;autoIncrement:false;comment:id" json:"id"`
	GroupID        string              `gorm:"size:64;index;not null;comment:聊天群id" json:"groupId"`
	UserID         string              `gorm:"size:64;index;not null;comment:成员id" json:"userId"`
	GroupNickName  string              `gorm:"size:128;comment:群昵称" json:"groupNickName"`
	GroupUserLevel int                 `gorm:"type:int;default:0;comment:群用户等级" json:"groupUserLevel"`
	MemberRole     string              `gorm:"size:128;comment:成员角色;default:'general'" json:"memberRole"`
	CreatedAt      localtime.LocalTime `gorm:"type:timestamp(3);not null;autoCreateTime;comment:创建时间" json:"createdAt"`
	UpdatedAt      localtime.LocalTime `gorm:"type:timestamp(3);not null;autoUpdateTime;comment:更新时间" json:"updatedAt"`
	DeletedAt      gorm.DeletedAt      `gorm:"index" json:"deletedAt"`

	Username    string `gorm:"->;-:migration" json:"username"`
	EmotionName string `gorm:"->;-:migration" json:"emotionName"`
	EmotionUrl  string `gorm:"->;-:migration" json:"emotionUrl"`
}

func (GroupMember) TableName() string {
	return "t_group_member"
}

func (GroupMember) TableComment() string {
	return "聊天群成员表"
}
