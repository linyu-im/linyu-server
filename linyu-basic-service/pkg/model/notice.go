package model

import (
	"encoding/json"

	"github.com/linyu-im/linyu-server/linyu-common/pkg/db"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/localtime"
	"gorm.io/gorm"
)

func init() {
	db.AddMigrateTable(&Notice{})
}

// Notice 通知表
type Notice struct {
	ID           string              `gorm:"size:64;primaryKey;autoIncrement:false;comment:通知id" json:"id"`
	UserID       string              `gorm:"size:64;index;not null;comment:接收通知的用户id" json:"userId"`
	SenderID     string              `gorm:"size:64;comment:通知发送方id，系统通知可为空system" json:"senderId"`
	Type         string              `gorm:"size:64;not null;comment:通知类型" json:"type"`
	NoticeSource string              `gorm:"size:64;comment:通知来源" json:"noticeSource"`
	Extra        json.RawMessage     `gorm:"type:json;comment:扩展数据" json:"extra"`
	CreatedAt    localtime.LocalTime `gorm:"type:timestamp(3);not null;autoCreateTime;comment:创建时间" json:"createdAt"`
	UpdatedAt    localtime.LocalTime `gorm:"type:timestamp(3);not null;autoUpdateTime;comment:更新时间" json:"updatedAt"`
	DeletedAt    gorm.DeletedAt      `gorm:"index" json:"deletedAt"`
}

type GroupNoticeExtra struct {
	GroupID     string `json:"groupId"`
	Status      string `json:"status"`
	LeaveUserID string `json:"leaveUserId"`
}

func (Notice) TableName() string {
	return "t_notice"
}

func (Notice) TableComment() string {
	return "通知表"
}
