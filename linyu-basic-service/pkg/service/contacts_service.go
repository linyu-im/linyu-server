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

func (s *contactsService) ContactsList(userId string) ([]*basicModel.Contacts, error) {
	list, err := basicDao.ContactsDao.ContactsList(db.MysqlDB, userId)
	return list, err
}

func (s *contactsService) ContactsRelDelete(userId string, param *basicParam.ContactsRelDeleteParam) error {
	contacts, err := basicDao.ContactsDao.GetById(db.MysqlDB, param.ContactsId)
	if contacts == nil || contacts.UserID != userId {
		return errors.New("param.error")
	}
	if contacts.Type == constant.ContactsType.User {
		err = db.MysqlDB.Transaction(func(tx *gorm.DB) error {
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
