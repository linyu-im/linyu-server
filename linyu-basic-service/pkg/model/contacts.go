package model

import (
	"github.com/linyu-im/linyu-server/linyu-common/pkg/db"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/localtime"
	"gorm.io/gorm"
)

func init() {
	db.AddMigrateTable(&Contacts{})
}

// Contacts 通讯录
type Contacts struct {
	ID        string              `gorm:"size:64;primaryKey;autoIncrement:false;comment:通讯录id" json:"id"`
	UserID    string              `gorm:"size:64;index;not null;comment:用户id" json:"userId"`
	PeerId    string              `gorm:"size:64;index;not null;comment:对方的id" json:"peerId"`
	Remark    string              `gorm:"size:64;comment:备注" json:"remark"`
	Tag       string              `gorm:"type:text;comment:标签" json:"tag"`
	IsBack    bool                `gorm:"default:0;comment:是否拉黑" json:"isBack"`
	IsTop     bool                `gorm:"default:0;comment:是否置顶" json:"isTop"`
	IsConcern bool                `gorm:"default:0;comment:是否关心" json:"isConcern"`
	IsMute    bool                `gorm:"default:0;comment:是否免打扰" json:"isMute"`
	PeerType  string              `gorm:"size:64;comment:对方的类型" json:"peerType"`
	Status    string              `gorm:"size:64;comment:状态" json:"status"`
	CreatedAt localtime.LocalTime `gorm:"type:timestamp(3);not null;autoCreateTime;comment:创建时间" json:"createdAt"`
	UpdatedAt localtime.LocalTime `gorm:"type:timestamp(3);not null;autoUpdateTime;comment:更新时间" json:"updatedAt"`
	DeletedAt gorm.DeletedAt      `gorm:"index" json:"deletedAt"`

	// 好友相关
	Username    string `gorm:"->;-:migration" json:"username"`
	UserLevel   int    `gorm:"->;-:migration" json:"userLevel"`
	EmotionName string `gorm:"->;-:migration" json:"emotionName"`
	EmotionUrl  string `gorm:"->;-:migration" json:"emotionUrl"`

	//群相关
	GroupName string `gorm:"->;-:migration" json:"groupName"`
	MemberNum int    `gorm:"->;-:migration" json:"memberNum"`
}

func (Contacts) TableName() string {
	return "t_contacts"
}

func (Contacts) TableComment() string {
	return "通讯录表"
}
