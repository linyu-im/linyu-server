package dao

import (
	"errors"

	basicModel "github.com/linyu-im/linyu-server/linyu-basic-service/pkg/model"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/localtime"
	"gorm.io/gorm"
)

var ChatDao = newChatDao()

func newChatDao() *chatDao {
	return &chatDao{}
}

type chatDao struct{}

func (d *chatDao) ChatList(db *gorm.DB, userId string) ([]*basicModel.Chat, error) {
	var chatList []*basicModel.Chat
	if err := db.Model(&basicModel.Chat{}).
		Where("user_id = ?", userId).
		Order("is_top DESC, updated_at DESC").
		Find(&chatList).Error; err != nil {
		return nil, err
	}
	return chatList, nil
}

func (d *chatDao) ContactsChatList(db *gorm.DB, userId string) ([]*basicModel.Chat, error) {
	var chatList []*basicModel.Chat

	err := db.Table("t_chat AS c").
		Select(`
			c.*,
			ct.remark AS peer_remark,
			ct.is_top AS peer_is_top,
			ct.is_mute AS peer_is_mute,
			CASE WHEN c.scene_type = 'user' THEN u.username ELSE g.name END AS peer_name,
			CASE WHEN c.scene_type = 'user' THEN u.avatar ELSE g.avatar END AS peer_avatar
		`).
		Joins("LEFT JOIN t_contacts ct ON ct.user_id = c.user_id AND ct.peer_id = c.peer_id AND ct.deleted_at IS NULL").
		Joins("LEFT JOIN t_user u ON u.id = c.peer_id AND u.deleted_at IS NULL").
		Joins("LEFT JOIN t_group g ON g.id = c.peer_id AND g.deleted_at IS NULL").
		Where("c.user_id = ? AND c.deleted_at IS NULL", userId).
		Order("ct.is_top DESC, c.updated_at DESC").
		Scan(&chatList).Error

	if err != nil {
		return nil, err
	}

	return chatList, nil
}

func (d *chatDao) GetChatByIdAndUserId(db *gorm.DB, userId string, chatId string) (*basicModel.Chat, error) {
	result := &basicModel.Chat{}
	if err := db.First(result, "id = ? AND user_id = ?", chatId, userId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return result, nil
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

func (d *chatDao) SetIsTopByIdAndUserId(db *gorm.DB, isTop bool, userId string, chatId string) error {
	chat, err := d.GetChatByIdAndUserId(db, userId, chatId)
	if err != nil {
		return err
	}
	if chat == nil {
		return errors.New("param.error")
	}
	if err = ContactsDao.SetIsTopByUserAndPeerId(db, isTop, userId, chat.PeerID); err != nil {
		return err
	}
	return d.touchUpdatedAtByIdAndUserId(db, userId, chatId)
}

func (d *chatDao) SetIsMuteByIdAndUserId(db *gorm.DB, isMute bool, userId string, chatId string) error {
	chat, err := d.GetChatByIdAndUserId(db, userId, chatId)
	if err != nil {
		return err
	}
	if chat == nil {
		return errors.New("param.error")
	}
	if err = ContactsDao.SetIsMuteByUserAndPeerId(db, isMute, userId, chat.PeerID); err != nil {
		return err
	}
	return d.touchUpdatedAtByIdAndUserId(db, userId, chatId)
}

func (d *chatDao) touchUpdatedAtByIdAndUserId(db *gorm.DB, userId string, chatId string) error {
	return db.Model(&basicModel.Chat{}).
		Where("id = ? AND user_id = ?", chatId, userId).
		Updates(map[string]interface{}{
			"updated_at": localtime.Now(),
		}).Error
}

func (d *chatDao) DeleteByIdAndUserId(db *gorm.DB, userId string, chatId string) error {
	return db.Delete(&basicModel.Chat{}, "user_id = ? AND id = ?", userId, chatId).Error
}

func (d *chatDao) ClearUnreadByIdAndUserId(db *gorm.DB, userId string, chatId string) error {
	return db.Model(&basicModel.Chat{}).
		Where("id = ? AND user_id = ?", chatId, userId).
		UpdateColumn("unread_num", 0).Error
}
