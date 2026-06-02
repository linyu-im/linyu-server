package model

import (
	"github.com/linyu-im/linyu-server/linyu-common/pkg/db"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/localtime"
	"gorm.io/gorm"
)

func init() {
	db.AddMigrateTable(&Enterprise{})
}

// Enterprise 企业
type Enterprise struct {
	ID               string              `gorm:"size:64;primaryKey;autoIncrement:false;comment:企业id" json:"id"`
	CreatorUserID    string              `gorm:"size:64;not null;comment:创建用户id" json:"creatorUserId"`
	EnterpriseNumber string              `gorm:"size:64;not null;comment:企业号" json:"enterpriseNumber"`
	Location         string              `gorm:"size:1024;comment:企业位置" json:"location"`
	Name             string              `gorm:"size:128;not null;comment:企业名称" json:"name"`
	Avatar           string              `gorm:"size:512;comment:企业头像URL" json:"avatar"`
	Describe         string              `gorm:"type:text;comment:企业描述" json:"describe"`
	OwnerUserID      string              `gorm:"size:64;not null;comment:企业主用户id" json:"ownerUserId"`
	EnterpriseTag    string              `gorm:"size:1024;comment:企业标签" json:"enterpriseTag"`
	MemberNum        int                 `gorm:"default:0;comment:企业成员数" json:"memberNum"`
	CreatedAt        localtime.LocalTime `gorm:"type:timestamp(3);not null;autoCreateTime;comment:创建时间" json:"createdAt"`
	UpdatedAt        localtime.LocalTime `gorm:"type:timestamp(3);not null;autoUpdateTime;comment:更新时间" json:"updatedAt"`
	DeletedAt        gorm.DeletedAt      `gorm:"index" json:"deletedAt"`

	Roles                 string              `gorm:"->;-:migration" json:"roles"`
	UserEnterpriseMembers []*EnterpriseMember `gorm:"-" json:"userEnterpriseMembers"`
}

func (Enterprise) TableName() string {
	return "t_enterprise"
}

func (Enterprise) TableComment() string {
	return "企业表"
}
