package service

import (
	"errors"
	basicDao "github.com/linyu-im/linyu-server/linyu-basic-service/internal/dao"
	basicModel "github.com/linyu-im/linyu-server/linyu-basic-service/pkg/model"
	basicParam "github.com/linyu-im/linyu-server/linyu-basic-service/pkg/param"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/constant"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/db"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/event"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/event/eventbus"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/utils"
	"gorm.io/gorm"
)

var ApplyService = newApplyService()

func newApplyService() *applyService {
	return &applyService{}
}

type applyService struct{}

func (s applyService) ApplyAddContacts(userId string, param *basicParam.ApplyAddContactsParam) error {
	//验证是否已经添加
	is := basicDao.ContactsDao.IsContactByUserAndPeer(db.MysqlDB, userId, param.PeerId)
	if is {
		return errors.New("basic.contacts.rel-already-exists")
	}
	apply := &basicModel.Apply{
		ID:       utils.GenerateSfIDString(),
		UserID:   userId,
		PeerID:   param.PeerId,
		Describe: param.Describe,
		Type:     constant.ApplyType.AddContacts,
		Status:   constant.ApplyStatus.Wait,
	}
	err := basicDao.ApplyDao.Create(db.MysqlDB, apply)
	if err != nil {
		return err
	}
	//发送申请消息
	_ = eventbus.DefaultEventBus.Publish(event.WsDataEvent{
		FromUserId: userId,
		ToUserIds:  []string{param.PeerId},
		Type:       constant.WsDataType.Apply,
		Content:    apply,
	})
	return nil
}

func (s applyService) ApplyAgreeContacts(userId string, param *basicParam.ApplyAgreeContactsParam) error {
	apply := basicDao.ApplyDao.GetById(db.MysqlDB, param.ApplyId)
	if apply == nil {
		return errors.New("common.data-not-exist")
	}
	if apply.PeerID != userId {
		return errors.New("param.error")
	}
	// 开始事务
	err := db.MysqlDB.Transaction(func(tx *gorm.DB) error {
		// 更新申请信息
		apply.Status = constant.ApplyStatus.Agree
		if err := basicDao.ApplyDao.Update(tx, apply); err != nil {
			return err
		}
		// 新增通讯双方关系
		if err := createContactIfNotExist(tx, userId, apply.UserID); err != nil {
			return err
		}
		if err := createContactIfNotExist(tx, apply.UserID, userId); err != nil {
			return err
		}
		return nil
	})
	return err
}

func (s applyService) ApplyReject(userId string, param *basicParam.ApplyRejectParam) error {
	apply := basicDao.ApplyDao.GetById(db.MysqlDB, param.ApplyId)
	if apply == nil {
		return errors.New("common.data-not-exist")
	}
	if apply.PeerID != userId || apply.Status != constant.ApplyStatus.Wait {
		return errors.New("param.error")
	}
	// 更新申请信息
	apply.Status = constant.ApplyStatus.Reject
	if err := basicDao.ApplyDao.Update(db.MysqlDB, apply); err != nil {
		return err
	}
	return nil
}

func (s applyService) ApplyCancel(userId string, param *basicParam.ApplyCancelParam) error {
	apply := basicDao.ApplyDao.GetById(db.MysqlDB, param.ApplyId)
	if apply == nil {
		return errors.New("common.data-not-exist")
	}
	if apply.UserID != userId || apply.Status != constant.ApplyStatus.Wait {
		return errors.New("param.error")
	}
	// 更新申请信息
	apply.Status = constant.ApplyStatus.Cancel
	if err := basicDao.ApplyDao.Update(db.MysqlDB, apply); err != nil {
		return err
	}
	return nil
}

func (s applyService) ApplyList(userId string) ([]*basicModel.Apply, error) {
	return basicDao.ApplyDao.ApplyListAndPeer(db.MysqlDB, userId)
}

func createContactIfNotExist(tx *gorm.DB, userID, peerID string) error {
	if !basicDao.ContactsDao.IsContactByUserAndPeer(tx, userID, peerID) {
		return basicDao.ContactsDao.Create(tx, &basicModel.Contacts{
			ID:     utils.GenerateSfIDString(),
			UserID: userID,
			PeerId: peerID,
			Type:   constant.ContactsType.User,
		})
	}
	return nil
}
