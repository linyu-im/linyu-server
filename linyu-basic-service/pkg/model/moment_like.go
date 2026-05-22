package model

import (
	"github.com/linyu-im/linyu-server/linyu-common/pkg/db"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/localtime"
)

func init() {
	db.AddMigrateTable(&MomentLike{})
}

// MomentLike 过往点赞表
type MomentLike struct {
	ID        string              `gorm:"size:64;primaryKey;autoIncrement:false" json:"id"`
	MomentID  string              `gorm:"size:64;index;comment:朋友圈id" json:"momentId"`
	UserID    string              `gorm:"size:64;index;comment:用户id" json:"userId"`
	CreatedAt localtime.LocalTime `gorm:"type:timestamp(3);not null;autoCreateTime;comment:创建时间" json:"createdAt"`

	Username string `gorm:"->;-:migration" json:"username"`
}

func (MomentLike) TableName() string {
	return "t_moment_like"
}

func (MomentLike) TableComment() string {
	return "过往点赞表"
}
