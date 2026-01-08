package model

import (
	"github.com/linyu-im/linyu-server/linyu-common/pkg/db"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/localtime"
	"gorm.io/gorm"
)

func init() {
	db.MysqlAddMigrateTable(&Message{})
}

// Message 消息表
type Message struct {
	ID         string              `gorm:"size:64;primaryKey;autoIncrement:false;comment:消息id" json:"id"`
	FromID     string              `gorm:"size:64;not null;comment:消息发送方id" json:"fromId"`
	ToID       string              `gorm:"size:64;not null;comment:消息接受方id" json:"toId"`
	Type       string              `gorm:"size:64;comment:消息类型" json:"type"`
	IsShowTime bool                `gorm:"comment:是否显示时间;default:0" json:"isShowTime"`
	Content    string              `gorm:"type:text;comment:消息内容" json:"content"`
	Status     string              `gorm:"size:500;comment:消息状态" json:"status"`
	Source     string              `gorm:"size:64;not null;comment:消息源" json:"source"`
	QuoteMsgId string              `gorm:"size:64;comment:引用消息的id" json:"quoteMsgId"`
	CreatedAt  localtime.LocalTime `gorm:"type:timestamp(3);not null;autoCreateTime;comment:创建时间" json:"createdAt"`
	UpdatedAt  localtime.LocalTime `gorm:"type:timestamp(3);not null;autoUpdateTime;comment:更新时间" json:"updatedAt"`
	DeletedAt  gorm.DeletedAt      `gorm:"index" json:"deletedAt"`
}

func (Message) TableName() string {
	return "t_message"
}

func (Message) TableComment() string {
	return "消息表"
}
