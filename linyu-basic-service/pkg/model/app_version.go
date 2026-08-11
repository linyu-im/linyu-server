package model

import (
	"github.com/linyu-im/linyu-server/linyu-common/pkg/db"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/localtime"
	"gorm.io/gorm"
)

func init() {
	db.AddMigrateTable(&AppVersion{})
}

// AppVersion 客户端版本配置
type AppVersion struct {
	ID                    string              `gorm:"size:64;primaryKey;autoIncrement:false;comment:id" json:"id"`
	Platform              string              `gorm:"size:32;uniqueIndex;not null;comment:平台" json:"platform"`
	LatestVersion         string              `gorm:"size:64;comment:最新版本号" json:"latestVersion"`
	LatestVersionCode     int                 `gorm:"type:int;default:0;comment:最新版本号(int)" json:"latestVersionCode"`
	MinSupportVersion     string              `gorm:"size:64;comment:最小支持版本号" json:"minSupportVersion"`
	MinSupportVersionCode int                 `gorm:"type:int;default:0;comment:最小支持版本号(int)" json:"minSupportVersionCode"`
	DownloadUrl           string              `gorm:"size:512;comment:下载地址" json:"downloadUrl"`
	UpdateDesc            string              `gorm:"type:text;comment:更新说明" json:"updateDesc"`
	Enabled               bool                `gorm:"default:1;comment:是否启用" json:"enabled"`
	CreatedAt             localtime.LocalTime `gorm:"type:timestamp(3);not null;autoCreateTime;comment:创建时间" json:"createdAt"`
	UpdatedAt             localtime.LocalTime `gorm:"type:timestamp(3);not null;autoUpdateTime;comment:更新时间" json:"updatedAt"`
	DeletedAt             gorm.DeletedAt      `gorm:"index" json:"deletedAt"`
}

func (AppVersion) TableName() string {
	return "t_app_version"
}

func (AppVersion) TableComment() string {
	return "客户端版本配置表"
}
