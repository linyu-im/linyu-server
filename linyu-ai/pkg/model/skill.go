package model

import (
	"github.com/linyu-im/linyu-server/linyu-common/pkg/db"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/localtime"
	"gorm.io/gorm"
)

func init() {
	db.AddMigrateTable(&Skill{})
}

// Skill Skill表
type Skill struct {
	ID           string              `gorm:"size:64;primaryKey;autoIncrement:false;comment:id" json:"id"`
	Name         string              `gorm:"size:512;not null;comment:名称" json:"name"`
	Description  string              `gorm:"type:text;comment:描述" json:"description"`
	Category     string              `gorm:"size:128;index;comment:分类" json:"category"`
	Version      string              `gorm:"size:64;comment:版本号" json:"version"`
	Author       string              `gorm:"size:512;comment:作者" json:"author"`
	IconUrl      string              `gorm:"size:1024;comment:图标地址" json:"iconUrl"`
	Featured     bool                `gorm:"default:0;index;comment:是否精选" json:"featured"`
	Capabilities []string            `gorm:"type:text;serializer:json;comment:能力标签" json:"capabilities"`
	Content      string              `gorm:"type:longtext;comment:skill内容" json:"content"`
	CreatedAt    localtime.LocalTime `gorm:"type:timestamp(3);not null;autoCreateTime;comment:创建时间" json:"createdAt"`
	UpdatedAt    localtime.LocalTime `gorm:"type:timestamp(3);not null;autoUpdateTime;comment:更新时间" json:"updatedAt"`
	DeletedAt    gorm.DeletedAt      `gorm:"index" json:"deletedAt"`
}

func (Skill) TableName() string {
	return "t_skill"
}

func (Skill) TableComment() string {
	return "Skill表"
}
