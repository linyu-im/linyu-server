package model

import (
	"github.com/linyu-im/linyu-server/linyu-common/pkg/db"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/localtime"
)

func init() {
	db.AddMigrateTable(&MomentSetting{})
}

// MomentSetting 过往设置表
type MomentSetting struct {
	UserID     string              `gorm:"size:64;primaryKey;comment:用户id" json:"userId"`
	BgUrl      string              `gorm:"size:512;comment:过往背景URL" json:"bgUrl"`
	ExpireDays int64               `gorm:"type:bigint;comment:朋友圈可见天数 0=永久 -1=不可见" json:"expireDays"`
	CreatedAt  localtime.LocalTime `gorm:"type:timestamp(3);autoCreateTime" json:"createdAt"`
	UpdatedAt  localtime.LocalTime `gorm:"type:timestamp(3);autoUpdateTime" json:"updatedAt"`
}

func (MomentSetting) TableName() string {
	return "t_moment_setting"
}

func (MomentSetting) TableComment() string {
	return "过往设置表"
}
