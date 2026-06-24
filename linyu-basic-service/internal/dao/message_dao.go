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

func (r messageDao) ListMessagesBySessionIDSinceMsgID(db *gorm.DB, sessionID string, sinceMsgId string) ([]*model.Message, error) {
	var messages []*model.Message
	query := db.Where("session_id = ?", sessionID)
	if sinceMsgId != "" {
		query = query.Where("id > ?", sinceMsgId)
	}
	if err := query.Order("created_at ASC").Find(&messages).Error; err != nil {
		return nil, err
	}
	return messages, nil
}

func (r messageDao) PageMessagesBySessionID(db *gorm.DB, sessionID string, page int, pageSize int) ([]*model.Message, int64, error) {

	var messages []*model.Message
	var total int64

	if page <= 0 {
		page = 1
	}

	if pageSize <= 0 {
		pageSize = 10
	}

	baseQuery := db.Model(&model.Message{}).
		Where("session_id = ?", sessionID)

	if err := baseQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize

	if err := db.Model(&model.Message{}).
		Where("session_id = ?", sessionID).
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&messages).Error; err != nil {

		return nil, 0, err
	}

	return messages, total, nil
}
