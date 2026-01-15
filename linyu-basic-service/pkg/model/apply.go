package model

import (
	"github.com/linyu-im/linyu-server/linyu-common/pkg/db"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/localtime"
	"gorm.io/gorm"
)

func init() {
	db.MysqlAddMigrateTable(&Apply{})
}

// Apply 申请相关表
type Apply struct {
	ID        string              `gorm:"size:64;primaryKey;autoIncrement:false;comment:id" json:"id"`
	UserID    string              `gorm:"size:64;not null;comment:用户id" json:"userId"`
	PeerID    string              `gorm:"size:64;comment:对方的id" json:"peerId"`
	Type      string              `gorm:"size:64;default:null;comment:类型" json:"type"`
	Describe  string              `gorm:"type:text;default:null;comment:描述" json:"describe"`
	Status    string              `gorm:"size:64;comment:状态" json:"status"`
	CreatedAt localtime.LocalTime `gorm:"type:timestamp(3);not null;autoCreateTime;comment:创建时间" json:"createdAt"`
	UpdatedAt localtime.LocalTime `gorm:"type:timestamp(3);not null;autoUpdateTime;comment:更新时间" json:"updatedAt"`
	DeletedAt gorm.DeletedAt      `gorm:"index" json:"deletedAt"`
}

func (Apply) TableName() string {
	return "t_apply"
}

func (Apply) TableComment() string {
	return "申请相关表"
}
