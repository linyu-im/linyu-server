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
