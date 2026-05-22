package model

import (
	"github.com/linyu-im/linyu-server/linyu-common/pkg/db"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/localtime"
)

func init() {
	db.AddMigrateTable(&MomentComment{})
}

// MomentComment 过往评论表
type MomentComment struct {
	ID            string              `gorm:"size:64;primaryKey;autoIncrement:false" json:"id"`
	MomentID      string              `gorm:"size:64;index;comment:过往id" json:"momentId"`
	UserID        string              `gorm:"size:64;index;comment:评论用户id" json:"userId"`
	ReplyUserID   string              `gorm:"size:64;comment:回复用户id" json:"replyUserId"`
	ReplyUsername string              `gorm:"size:64;comment:回复用户名称" json:"replyUsername"`
	ParentID      string              `gorm:"size:64;index;comment:父评论id" json:"parentId"`
	Content       string              `gorm:"type:text;comment:评论内容" json:"content"`
	CreatedAt     localtime.LocalTime `gorm:"type:timestamp(3);not null;autoCreateTime;comment:创建时间" json:"createdAt"`

	Username string `gorm:"->;-:migration" json:"username"`
}

func (MomentComment) TableName() string {
	return "t_moment_comment"
}

func (MomentComment) TableComment() string {
	return "过往评论表"
}
