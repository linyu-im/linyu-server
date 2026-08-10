package service

import (
	"errors"

	basicDao "github.com/linyu-im/linyu-server/linyu-basic-service/internal/dao"
	basicModel "github.com/linyu-im/linyu-server/linyu-basic-service/pkg/model"
	basicParam "github.com/linyu-im/linyu-server/linyu-basic-service/pkg/param"
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

func (s *contactsService) ContactsEnterpriseList(userId string) ([]*basicModel.Contacts, error) {
	list, err := basicDao.ContactsDao.ContactsEnterpriseList(db.RDB, userId)
	return list, err
}

func (s *contactsService) ContactsSearch(userId string, param *basicParam.ContactsSearchParam) (*basicParam.ContactsSearchResult, error) {
	friends, err := basicDao.ContactsDao.ContactsFriendSearch(db.RDB, userId, param.Keyword)
	if err != nil {
		return nil, err
	}
	groups, err := basicDao.ContactsDao.ContactsGroupSearch(db.RDB, userId, param.Keyword)
	if err != nil {
		return nil, err
	}
	if friends == nil {
		friends = []*basicModel.Contacts{}
	}
	if groups == nil {
		groups = []*basicModel.Contacts{}
	}
	return &basicParam.ContactsSearchResult{
		Friends: friends,
		Groups:  groups,
	}, nil
}

func (s *contactsService) IsFriend(userId string, peerId string) bool {
	return basicDao.ContactsDao.IsFriend(db.RDB, userId, peerId)
}

func (s *contactsService) IsFriendBothOr(userId string, peerId string) bool {
	return basicDao.ContactsDao.IsFriendBothOr(db.RDB, userId, peerId)
}

func (s *contactsService) IsFriendBothAnd(userId string, peerId string) bool {
	return basicDao.ContactsDao.IsFriendBothAnd(db.RDB, userId, peerId)
}

func (s *contactsService) UpdateRemark(userId string, param *basicParam.ContactsUpdateRemarkParam) error {
	return basicDao.ContactsDao.UpdateRemarkByUserAndPeerId(db.RDB, userId, param.PeerId, param.Remark)
}

func (s *contactsService) UpdateTag(userId string, param *basicParam.ContactsUpdateTagParam) error {
	return basicDao.ContactsDao.UpdateTagByUserAndPeerId(db.RDB, userId, param.PeerId, param.Tag)
}

func (s *contactsService) ContactsFriendDelete(currentUserId string, param *basicParam.ContactsFriendDeleteParam) error {
	if !s.IsFriend(currentUserId, param.UserId) {
		return errors.New("param.error")
	}
	err := db.RDB.Transaction(func(tx *gorm.DB) error {
		// 删除通讯录关系
		if err := basicDao.ContactsDao.UnscopedDeleteByUserAndPeerId(tx, currentUserId, param.UserId); err != nil {
			return err
		}
		// 删除聊天列表
		if err := basicDao.ChatDao.UnscopedDeleteByUserIdAndPeerId(tx, currentUserId, param.UserId); err != nil {
			return err
		}
		return nil
	})
	return err
}
