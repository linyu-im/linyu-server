package model

import (
	"github.com/linyu-im/linyu-server/linyu-common/pkg/db"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/localtime"
	"gorm.io/gorm"
)

func init() {
	db.MysqlAddMigrateTable(&Contacts{})
}

// Contacts 通讯录
type Contacts struct {
	ID        string              `gorm:"size:64;primaryKey;autoIncrement:false;comment:好友表id" json:"id"`
	UserID    string              `gorm:"size:64;index;not null;comment:用户id" json:"userId"`
	PeerId    string              `gorm:"size:64;not null;comment:对方的id" json:"peerId"`
	Remark    string              `gorm:"size:64;comment:备注" json:"remark"`
	IsBack    bool                `gorm:"comment:是否拉黑;default:0" json:"isBack"`
	IsConcern bool                `gorm:"comment:是否关心;default:0" json:"isConcern"`
	Type      string              `gorm:"size:64;comment:类型" json:"type"`
	Status    string              `gorm:"size:64;comment:状态" json:"status"`
	CreatedAt localtime.LocalTime `gorm:"type:timestamp(3);not null;autoCreateTime;comment:创建时间" json:"createdAt"`
	UpdatedAt localtime.LocalTime `gorm:"type:timestamp(3);not null;autoUpdateTime;comment:更新时间" json:"updatedAt"`
	DeletedAt gorm.DeletedAt      `gorm:"index" json:"deletedAt"`
}

func (Contacts) TableName() string {
	return "t_contacts"
}

func (Contacts) TableComment() string {
	return "通讯录表"
}
