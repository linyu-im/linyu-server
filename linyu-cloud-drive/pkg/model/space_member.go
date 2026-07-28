package model

import (
	"github.com/linyu-im/linyu-server/linyu-common/pkg/db"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/localtime"
	"gorm.io/gorm"
)

func init() {
	db.AddMigrateTable(&SpaceMember{})
}

// SpaceMember 空间成员表
type SpaceMember struct {
	ID         string              `gorm:"size:64;primaryKey;autoIncrement:false;comment:id" json:"id"`
	SpaceID    string              `gorm:"size:64;not null;uniqueIndex:uk_space_member,priority:1;comment:空间id" json:"spaceId"`
	UserID     string              `gorm:"size:64;not null;uniqueIndex:uk_space_member,priority:2;comment:成员用户id" json:"userId"`
	MemberRole string              `gorm:"size:32;not null;default:member;comment:成员角色 owner/admin/member" json:"memberRole"`
	Status     string              `gorm:"size:32;index;not null;default:active;comment:状态 active/disabled" json:"status"`
	CreatedAt  localtime.LocalTime `gorm:"type:timestamp(3);autoCreateTime;comment:创建时间" json:"createdAt"`
	UpdatedAt  localtime.LocalTime `gorm:"type:timestamp(3);autoUpdateTime;comment:更新时间" json:"updatedAt"`
	DeletedAt  gorm.DeletedAt      `gorm:"index" json:"deletedAt"`
}

func (SpaceMember) TableName() string {
	return "t_space_member"
}

func (SpaceMember) TableComment() string {
	return "空间成员表"
}
