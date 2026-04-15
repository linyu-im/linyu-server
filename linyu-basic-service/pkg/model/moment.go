package model

import (
	"github.com/linyu-im/linyu-server/linyu-common/pkg/db"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/localtime"
	"gorm.io/gorm"
)

func init() {
	db.AddMigrateTable(&Moment{})
}

type MediaContent struct {
	Url       string `json:"url"`
	ThumbURL  string `json:"thumbUrl"`
	MediaType string `json:"mediaType"`
	Sort      int    `json:"sort"`
}

// Moment 过往表
type Moment struct {
	ID            string              `gorm:"size:64;primaryKey;autoIncrement:false;comment:id" json:"id"`
	UserID        string              `gorm:"size:64;not null;index:idx_moment_user_deleted_time;comment:发布者用户id" json:"userId"`
	TextContent   string              `gorm:"type:text;comment:文字内容" json:"textContent"`
	MediaContents []*MediaContent     `gorm:"type:text;serializer:json;comment:媒体内容" json:"mediaType"`
	VisibleType   string              `gorm:"size:32;index:idx_visible_type;comment:可见性类型:all/private/include/exclude" json:"VisibleType"`
	Location      string              `gorm:"size:512;comment:位置信息" json:"location"`
	CreatedAt     localtime.LocalTime `gorm:"type:timestamp(3);not null;autoCreateTime;index:idx_moment_user_deleted_time,sort:desc;comment:创建时间" json:"createdAt"`
	UpdatedAt     localtime.LocalTime `gorm:"type:timestamp(3);not null;autoUpdateTime;comment:更新时间" json:"updatedAt"`
	DeletedAt     gorm.DeletedAt      `gorm:"index:idx_moment_user_deleted_time;comment:删除时间" json:"deletedAt"`
}

func (Moment) TableName() string {
	return "t_moment"
}

func (Moment) TableComment() string {
	return "过往表"
}
