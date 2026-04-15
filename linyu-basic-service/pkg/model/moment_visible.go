package model

import (
	"github.com/linyu-im/linyu-server/linyu-common/pkg/db"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/localtime"
)

func init() {
	db.AddMigrateTable(&MomentVisible{})
}

// MomentVisible 过往可见权限表
type MomentVisible struct {
	ID          string              `gorm:"size:64;primaryKey;autoIncrement:false;comment:id" json:"id"`
	MomentID    string              `gorm:"size:64;index;comment:过往id" json:"momentId"`
	UserID      string              `gorm:"size:64;index;comment:用户id" json:"userId"`
	VisibleType string              `gorm:"size:32;comment:include/exclude" json:"visibleType"`
	CreatedAt   localtime.LocalTime `gorm:"type:timestamp(3);not null;autoCreateTime;comment:创建时间" json:"createdAt"`
}

func (MomentVisible) TableName() string {
	return "t_moment_visible"
}

func (MomentVisible) TableComment() string {
	return "过往可见权限表"
}
