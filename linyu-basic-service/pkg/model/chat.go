package model

import (
	"github.com/linyu-im/linyu-server/linyu-common/pkg/db"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/localtime"
	"gorm.io/gorm"
)

func init() {
	db.AddMigrateTable(&Chat{})
}

type Chat struct {
	ID             string              `gorm:"size:64;primaryKey;autoIncrement:false;comment:id" json:"id"`
	UserID         string              `gorm:"size:64;index;not null;comment:用户id" json:"userId"`
	PeerID         string              `gorm:"size:64;index;not null;comment:会话对方id" json:"peerId"`
	UnreadNum      int                 `gorm:"type:int;default:0;comment:未读消息数量" json:"unreadNum"`
	LastMsgContent any                 `gorm:"type:text;serializer:json;comment:最后消息内容" json:"lastMsgContent"`
	Type           string              `gorm:"size:64;comment:类型" json:"type"`
	Status         string              `gorm:"size:64;comment:状态" json:"status"`
	CreatedAt      localtime.LocalTime `gorm:"type:timestamp(3);not null;autoCreateTime;comment:创建时间" json:"createdAt"`
	UpdatedAt      localtime.LocalTime `gorm:"type:timestamp(3);not null;autoUpdateTime;comment:更新时间" json:"updatedAt"`
	DeletedAt      gorm.DeletedAt      `gorm:"index" json:"deletedAt"`

	PeerName   string `gorm:"->;-:migration" json:"peerName"`
	PeerAvatar string `gorm:"->;-:migration" json:"peerAvatar"`
	PeerRemark string `gorm:"->;-:migration" json:"peerRemark"`
	PeerIsTop  bool   `gorm:"->;-:migration" json:"peerIsTop"`
	PeerIsMute bool   `gorm:"->;-:migration" json:"peerIsMute"`
}

func (Chat) TableName() string {
	return "t_chat"
}

func (Chat) TableComment() string {
	return "聊天列表"
}
