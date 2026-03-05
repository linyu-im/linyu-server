package dao

import (
	"github.com/linyu-im/linyu-server/linyu-basic-service/pkg/model"
	"gorm.io/gorm"
)

var MessageDao = newMessageDao()

func newMessageDao() *messageDao {
	return &messageDao{}
}

type messageDao struct{}

func (r messageDao) Create(db *gorm.DB, message *model.Message) error {
	if err := db.Create(message).Error; err != nil {
		return err
	}
	return nil
}

func (r messageDao) GetLatestMessagesBySessionID(db *gorm.DB, sessionID string, limit int) []*model.Message {
	var messages []*model.Message
	result := db.Where("session_id = ?", sessionID).Order("created_at DESC").
		Limit(limit).Find(&messages)

	if result.Error != nil {
		return []*model.Message{}
	}

	return messages
}
