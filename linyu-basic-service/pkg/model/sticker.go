package model

import (
	"github.com/linyu-im/linyu-server/linyu-common/pkg/db"
	initsticker "github.com/linyu-im/linyu-server/linyu-common/pkg/init/sticker"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/localtime"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func init() {
	db.AddMigrateTable(&Sticker{})
}

// Sticker 表情表
type Sticker struct {
	ID            string              `gorm:"size:64;primaryKey;autoIncrement:false;comment:id" json:"id"`
	Name          string              `gorm:"size:128;not null;comment:表情名称" json:"name"`
	IconUrl       string              `gorm:"size:512;comment:表情图标地址" json:"iconUrl"`
	Type          string              `gorm:"size:64;comment:类型" json:"type"`
	IconValue     string              `gorm:"size:512;comment:表情值" json:"iconValue"`
	StickerPackID string              `gorm:"size:64;index;not null;comment:表情分组id" json:"stickerPackId"`
	CreatedAt     localtime.LocalTime `gorm:"type:timestamp(3);not null;autoCreateTime;comment:创建时间" json:"createdAt"`
	UpdatedAt     localtime.LocalTime `gorm:"type:timestamp(3);not null;autoUpdateTime;comment:更新时间" json:"updatedAt"`
	DeletedAt     gorm.DeletedAt      `gorm:"index" json:"deletedAt"`
}

func (Sticker) TableName() string {
	return "t_sticker"
}

func (Sticker) TableComment() string {
	return "表情表"
}

func (Sticker) TableInit(db *gorm.DB) error {
	for _, item := range initsticker.All() {
		err := db.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			DoNothing: true,
		}).Create(&Sticker{
			ID:            item.ID,
			Name:          item.Name,
			IconUrl:       item.IconUrl,
			Type:          item.Type,
			IconValue:     item.IconValue,
			StickerPackID: item.StickerPackID,
		}).Error
		if err != nil {
			return err
		}
	}
	return nil
}
