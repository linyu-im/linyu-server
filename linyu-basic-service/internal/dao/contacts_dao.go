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

func (d *contactsDao) ContactsList(db *gorm.DB, userId string) ([]*basicModel.Contacts, error) {
	var list []*basicModel.Contacts
	if err := db.Where("user_id = ?", userId).Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (d *contactsDao) GetById(db *gorm.DB, contactsId string) (*basicModel.Contacts, error) {
	result := &basicModel.Contacts{}
	if err := db.First(result, "id = ?", contactsId).Error; err != nil {
		return nil, err
	}
	return result, nil
}

func (d *contactsDao) UnscopedDeleteByUserAndPeerId(db *gorm.DB, userId string, peerId string) error {
	result := db.Unscoped().Where("user_id = ? AND peer_id = ?", userId, peerId).Delete(&basicModel.Contacts{}).Error
	return result
}
