package model

import (
	"github.com/linyu-im/linyu-server/linyu-common/pkg/db"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/localtime"
	"gorm.io/gorm"
)

func init() {
	db.MysqlAddMigrateTable(&Chat{})
}

type Chat struct {
	ID             string              `gorm:"size:64;primaryKey;autoIncrement:false;comment:id" json:"id"`
	UserID         string              `gorm:"size:64;not null;index;comment:用户id" json:"userId"`
	PeerID         string              `gorm:"size:64;not null;index:idx_userid_peerid;comment:会话对方id" json:"peerId	"`
	IsTop          bool                `gorm:"default:0;comment:是否置顶" json:"isTop"`
	UnreadNum      int                 `gorm:"type:int;default:0;comment:未读消息数量" json:"unreadNum"`
	LastMsgContent *Message            `gorm:"type:text;serializer:json;comment:最后消息内容" json:"lastMsgContent"`
	Type           string              `gorm:"size:64;comment:类型" json:"type"`
	Status         string              `gorm:"size:64;comment:状态" json:"status"`
	CreatedAt      localtime.LocalTime `gorm:"type:timestamp(3);not null;autoCreateTime;comment:创建时间" json:"createdAt"`
	UpdatedAt      localtime.LocalTime `gorm:"type:timestamp(3);not null;autoUpdateTime;comment:更新时间" json:"updatedAt"`
	DeletedAt      gorm.DeletedAt      `gorm:"index" json:"deletedAt"`
}

func (Chat) TableName() string {
	return "t_chat"
}

func (Chat) TableComment() string {
	return "聊天列表"
}
