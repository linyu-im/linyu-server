package model

import (
	"github.com/linyu-im/linyu-server/linyu-common/pkg/db"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/localtime"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func init() {
	db.AddMigrateTable(&StickerPack{})
}

// StickerPack 表情分组表
type StickerPack struct {
	ID          string              `gorm:"size:64;primaryKey;autoIncrement:false;comment:id" json:"id"`
	Name        string              `gorm:"size:128;not null;comment:分组名称" json:"name"`
	Description string              `gorm:"type:text;comment:描述" json:"description"`
	PackIconUrl string              `gorm:"size:512;comment:分组图标地址" json:"packIconUrl"`
	CreatedAt   localtime.LocalTime `gorm:"type:timestamp(3);not null;autoCreateTime;comment:创建时间" json:"createdAt"`
	UpdatedAt   localtime.LocalTime `gorm:"type:timestamp(3);not null;autoUpdateTime;comment:更新时间" json:"updatedAt"`
	DeletedAt   gorm.DeletedAt      `gorm:"index" json:"deletedAt"`

	Stickers []*Sticker `gorm:"->;-:migration" json:"stickers"`
}

func (StickerPack) TableName() string {
	return "t_sticker_pack"
}

func (StickerPack) TableComment() string {
	return "表情分组表"
}

func (StickerPack) TableInit(db *gorm.DB) error {
	datas := []StickerPack{
		{
			ID:          "1",
			Name:        "米游兔",
			Description: "米游兔官方表情包",
			PackIconUrl: "https://upload-bbs.mihoyo.com/upload/2022/11/14/62158b26331a045bbaabf48d1fc5c8eb_980401343708652819.png",
		},
	}
	for _, item := range datas {
		err := db.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			DoNothing: true,
		}).Create(&item).Error
		if err != nil {
			return err
		}
	}
	return nil
}
