package dao

import (
	basicModel "github.com/linyu-im/linyu-server/linyu-basic-service/pkg/model"
	"gorm.io/gorm"
)

var ContactsDao = newContactsDao()

func newContactsDao() *contactsDao {
	return &contactsDao{}
}

type contactsDao struct{}

func (d *contactsDao) IsContactByUserAndPeer(db *gorm.DB, userId string, peerId string) bool {
	var count int64
	err := db.Model(&basicModel.Contacts{}).
		Where("user_id = ? AND peer_id = ?", userId, peerId).
		Count(&count).
		Error
	if err != nil {
		return false
	}
	return count > 0
}

func (d *contactsDao) Create(db *gorm.DB, contacts *basicModel.Contacts) error {
	if err := db.Create(contacts).Error; err != nil {
		return err
	}
	return nil
}
