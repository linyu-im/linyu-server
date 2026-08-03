package model

import (
	"github.com/linyu-im/linyu-server/linyu-common/pkg/db"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/localtime"
	"gorm.io/gorm"
)

func init() {
	db.AddMigrateTable(&Application{})
}

// Application 应用表
type Application struct {
	ID          string              `gorm:"size:64;primaryKey;autoIncrement:false;comment:id" json:"id"`
	AppName     string              `gorm:"size:512;not null;comment:应用名称" json:"appName"`
	Version     string              `gorm:"size:64;comment:版本号" json:"version"`
	Description string              `gorm:"type:text;comment:描述" json:"description"`
	AuthorID    string              `gorm:"size:64;index;comment:作者id" json:"authorId"`
	Author      string              `gorm:"size:512;comment:作者" json:"author"`
	Tags        []string            `gorm:"type:text;serializer:json;comment:标签" json:"tags"`
	AppType     string              `gorm:"size:64;index;comment:应用类型" json:"appType"`
	IconUrl     string              `gorm:"size:512;comment:图标地址" json:"iconUrl"`
	PluginUrl   string              `gorm:"size:512;comment:插件地址" json:"pluginUrl"`
	WebUrl      string              `gorm:"size:512;comment:web地址" json:"webUrl"`
	GetCount    int                 `gorm:"type:int;default:0;comment:获取次数" json:"getCount"`
	Score       float64             `gorm:"type:decimal(3,1);default:0;comment:评分" json:"score"`
	ScoreCount  int                 `gorm:"type:int;default:0;comment:评分人数" json:"scoreCount"`
	CreatedAt   localtime.LocalTime `gorm:"type:timestamp(3);not null;autoCreateTime;comment:创建时间" json:"createdAt"`
	UpdatedAt   localtime.LocalTime `gorm:"type:timestamp(3);not null;autoUpdateTime;comment:更新时间" json:"updatedAt"`
	DeletedAt   gorm.DeletedAt      `gorm:"index" json:"deletedAt"`
}

func (Application) TableName() string {
	return "t_application"
}

func (Application) TableComment() string {
	return "应用表"
}
