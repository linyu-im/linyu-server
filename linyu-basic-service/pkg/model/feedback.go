package model

import (
	"github.com/linyu-im/linyu-server/linyu-common/pkg/db"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/localtime"
	"gorm.io/gorm"
)

func init() {
	db.AddMigrateTable(&Feedback{})
}

// Feedback 反馈表
type Feedback struct {
	ID          string              `gorm:"size:64;primaryKey;autoIncrement:false;comment:id" json:"id"`
	UserID      string              `gorm:"size:64;index;not null;comment:用户id" json:"userId"`
	Title       string              `gorm:"size:256;not null;comment:标题" json:"title"`
	Description string              `gorm:"type:text;comment:补充描述" json:"description"`
	Images      []string            `gorm:"type:text;serializer:json;comment:反馈图片" json:"images"`
	CreatedAt   localtime.LocalTime `gorm:"type:timestamp(3);not null;autoCreateTime;comment:创建时间" json:"createdAt"`
	UpdatedAt   localtime.LocalTime `gorm:"type:timestamp(3);not null;autoUpdateTime;comment:更新时间" json:"updatedAt"`
	DeletedAt   gorm.DeletedAt      `gorm:"index" json:"deletedAt"`
}

func (Feedback) TableName() string {
	return "t_feedback"
}

func (Feedback) TableComment() string {
	return "反馈表"
}
