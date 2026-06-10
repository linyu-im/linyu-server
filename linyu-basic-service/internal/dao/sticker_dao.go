package dao

import (
	basicModel "github.com/linyu-im/linyu-server/linyu-basic-service/pkg/model"
	"gorm.io/gorm"
)

const DefaultStickerPackID = "default"

var StickerDao = newStickerDao()

func newStickerDao() *stickerDao {
	return &stickerDao{}
}

type stickerDao struct{}

func (d *stickerDao) DefaultStickerList(db *gorm.DB) ([]*basicModel.Sticker, error) {
	var list []*basicModel.Sticker
	if err := db.Where("sticker_pack_id = ?", DefaultStickerPackID).
		Order("id ASC").
		Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (d *stickerDao) AllStickerPackList(db *gorm.DB) ([]*basicModel.StickerPack, error) {
	var list []*basicModel.StickerPack
	if err := db.Order("id ASC").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (d *stickerDao) AllStickerList(db *gorm.DB) ([]*basicModel.Sticker, error) {
	var list []*basicModel.Sticker
	if err := db.Order("sticker_pack_id ASC, id ASC").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (d *stickerDao) FavoriteStickerList(db *gorm.DB, userId string) ([]*basicModel.Sticker, error) {
	var list []*basicModel.Sticker
	if err := db.Table("t_user_sticker_favorite f").
		Select("s.*").
		Joins("INNER JOIN t_sticker s ON s.id = f.sticker_id AND s.deleted_at IS NULL").
		Where("f.user_id = ? AND f.deleted_at IS NULL", userId).
		Order("f.created_at DESC").
		Scan(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}
