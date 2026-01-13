package dao

import (
	"errors"
	basicModel "github.com/linyu-im/linyu-server/linyu-basic-service/pkg/model"
	"gorm.io/gorm"
)

var ChatDao = newChatDao()

func newChatDao() *chatDao {
	return &chatDao{}
}

type chatDao struct{}

func (d *chatDao) ChatList(db *gorm.DB, userId string) ([]*basicModel.Chat, error) {
	var chatList []*basicModel.Chat
	if err := db.Where("user_id = ?", userId).Find(&chatList).Error; err != nil {
		return nil, err
	}
	return chatList, nil
}

func (d *chatDao) create(db *gorm.DB, chat *basicModel.Chat) error {
	if err := db.Create(chat).Error; err != nil {
		return err
	}
	return nil
}

func (d *chatDao) GetChatByUserAndPeer(db *gorm.DB, userId string, peerId string) (*basicModel.Chat, error) {
	result := &basicModel.Chat{}
	if err := db.First(result, "user_id = ? AND peer_id = ?", userId, peerId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return result, nil
}

func (d *chatDao) Create(db *gorm.DB, chat *basicModel.Chat) error {
	if err := db.Create(chat).Error; err != nil {
		return err
	}
	return nil
}

func (d *chatDao) Update(db *gorm.DB, chat *basicModel.Chat) error {
	if err := db.Updates(chat).Error; err != nil {
		return err
	}
	return nil
}
