package service

import (
	"errors"

	basicDao "github.com/linyu-im/linyu-server/linyu-basic-service/internal/dao"
	basicModel "github.com/linyu-im/linyu-server/linyu-basic-service/pkg/model"
	basicParam "github.com/linyu-im/linyu-server/linyu-basic-service/pkg/param"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/constant"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/db"
	"gorm.io/gorm"
)

var ContactsService = newContactsService()

func newContactsService() *contactsService {
	return &contactsService{}
}

type contactsService struct{}

func (s *contactsService) ContactsFriendList(userId string) ([]*basicModel.Contacts, error) {
	list, err := basicDao.ContactsDao.ContactsFriendList(db.RDB, userId)
	return list, err
}

func (s *contactsService) ContactsGroupList(userId string) ([]*basicModel.Contacts, error) {
	list, err := basicDao.ContactsDao.ContactsGroupList(db.RDB, userId)
	return list, err
}

func (s *contactsService) IsContacts(userId string, peerId string) bool {
	return basicDao.ContactsDao.IsContactByUserAndPeer(db.RDB, userId, peerId)
}

func (s *contactsService) ContactsRelDelete(userId string, param *basicParam.ContactsRelDeleteParam) error {
	contacts, err := basicDao.ContactsDao.GetById(db.RDB, param.ContactsId)
	if contacts == nil || contacts.UserID != userId {
		return errors.New("param.error")
	}
	if contacts.PeerType == constant.ContactsPeerType.Friend {
		err = db.RDB.Transaction(func(tx *gorm.DB) error {
			// 双方关系删除
			if err := basicDao.ContactsDao.UnscopedDeleteByUserAndPeerId(tx, contacts.UserID, contacts.PeerId); err != nil {
				return err
			}
			if err := basicDao.ContactsDao.UnscopedDeleteByUserAndPeerId(tx, contacts.PeerId, contacts.UserID); err != nil {
				return err
			}
			return nil
		})
	}
	return err
}
