package service

import (
	basicDao "github.com/linyu-im/linyu-server/linyu-basic-service/internal/dao"
	basicModel "github.com/linyu-im/linyu-server/linyu-basic-service/pkg/model"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/db"
)

var StickerService = newStickerService()

func newStickerService() *stickerService {
	return &stickerService{}
}

type stickerService struct{}

func (s stickerService) DefaultStickerList() ([]*basicModel.Sticker, error) {
	return basicDao.StickerDao.DefaultStickerList(db.RDB)
}

func (s stickerService) FavoriteStickerList(userId string) ([]*basicModel.Sticker, error) {
	return basicDao.StickerDao.FavoriteStickerList(db.RDB, userId)
}

func (s stickerService) StickerPackList() ([]*basicModel.StickerPack, error) {
	packs, err := basicDao.StickerDao.AllStickerPackList(db.RDB)
	if err != nil {
		return nil, err
	}
	stickers, err := basicDao.StickerDao.AllStickerList(db.RDB)
	if err != nil {
		return nil, err
	}
	stickerMap := make(map[string][]*basicModel.Sticker)
	for _, sticker := range stickers {
		stickerMap[sticker.StickerPackID] = append(stickerMap[sticker.StickerPackID], sticker)
	}
	for _, pack := range packs {
		pack.Stickers = stickerMap[pack.ID]
		if pack.Stickers == nil {
			pack.Stickers = []*basicModel.Sticker{}
		}
	}
	return packs, nil
}
