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

func (s applyService) ApplyAddFriend(userId string, param *basicParam.ApplyAddFriendParam) error {
	if !constant.ApplySource.Validate(param.ApplySource) {
		return errors.New("param.error")
	}
	//验证是否已经添加(双方)
	is := basicDao.ContactsDao.IsFriendBothAnd(db.RDB, userId, param.PeerId)
	if is {
		return errors.New("basic.contacts.friend-already-exists")
	}
	apply := &basicModel.Apply{
		ID:          utils.GenerateSfIDString(),
		UserID:      userId,
		PeerID:      param.PeerId,
		Describe:    param.Describe,
		ApplySource: param.ApplySource,
		Type:        constant.ApplyType.Friend,
		Status:      constant.ApplyStatus.Wait,
	}
	err := basicDao.ApplyDao.Create(db.RDB, apply)
	if err != nil {
		return err
	}
	//发送申请消息
	_ = eventbus.GlobalBus.Publish(event.WsDataEvent{
		FromUserId: userId,
		ToUserIds:  []string{param.PeerId},
		Data: &event.WsData{
			SeqId:   apply.ID,
			Type:    constant.WsDataType.Apply,
			Content: apply,
		},
	})
	return nil
}

func (s applyService) ApplyAddGroup(userId string, param *basicParam.ApplyAddGroupParam) error {
	if !constant.ApplySource.Validate(param.ApplySource) {
		return errors.New("param.error")
	}
	if basicDao.GroupMemberDao.GetGroupMemberByGroupIdAndUserId(db.RDB, param.GroupId, userId) != nil {
		return errors.New("basic.group.member-already-exists")
	}
	apply := &basicModel.Apply{
		ID:          utils.GenerateSfIDString(),
		UserID:      userId,
		PeerID:      param.GroupId,
		Describe:    param.Describe,
		ApplySource: param.ApplySource,
		Type:        constant.ApplyType.Group,
		Status:      constant.ApplyStatus.Wait,
	}
	err := basicDao.ApplyDao.Create(db.RDB, apply)
	if err != nil {
		return err
	}
	// 通知群管理员
	adminIds := basicDao.GroupMemberDao.GetAdminUserIdsByGroupId(db.RDB, param.GroupId)
	if len(adminIds) > 0 {
		_ = eventbus.GlobalBus.Publish(event.WsDataEvent{
			FromUserId: userId,
			ToUserIds:  adminIds,
			Data: &event.WsData{
				SeqId:   apply.ID,
				Type:    constant.WsDataType.Apply,
				Content: apply,
			},
		})
	}
	return nil
}

func (s applyService) ApplyAgreeFriend(userId string, param *basicParam.ApplyAgreeFriendParam) error {
	apply := basicDao.ApplyDao.GetById(db.RDB, param.ApplyId)
	if apply == nil {
		return errors.New("common.data-not-exist")
	}
	if apply.PeerID != userId {
		return errors.New("param.error")
	}
	// 开始事务
	err := db.RDB.Transaction(func(tx *gorm.DB) error {
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
	apply := basicDao.ApplyDao.GetById(db.RDB, param.ApplyId)
	if apply == nil {
		return errors.New("common.data-not-exist")
	}
	if apply.PeerID != userId || apply.Status != constant.ApplyStatus.Wait {
		return errors.New("param.error")
	}
	// 更新申请信息
	apply.Status = constant.ApplyStatus.Reject
	if err := basicDao.ApplyDao.Update(db.RDB, apply); err != nil {
		return err
	}
	return nil
}

func (s applyService) ApplyCancel(userId string, param *basicParam.ApplyCancelParam) error {
	apply := basicDao.ApplyDao.GetById(db.RDB, param.ApplyId)
	if apply == nil {
		return errors.New("common.data-not-exist")
	}
	if apply.UserID != userId || apply.Status != constant.ApplyStatus.Wait {
		return errors.New("param.error")
	}
	// 更新申请信息
	apply.Status = constant.ApplyStatus.Cancel
	if err := basicDao.ApplyDao.Update(db.RDB, apply); err != nil {
		return err
	}
	return nil
}

func (s applyService) ApplyFriendList(userId string) ([]*basicModel.Apply, error) {
	list, err := basicDao.ApplyDao.ApplyFriendList(db.RDB, userId)
	if err != nil {
		return nil, err
	}
	lastReadId := "0"
	if len(list) > 0 {
		lastReadId = list[0].ID
	}
	_ = basicDao.UserBadgeDao.UpsertLastReadID(db.RDB, userId, constant.BadgeCode.NewFriend, lastReadId)
	return list, nil
}

func (s applyService) ApplyGroupList(userId string) ([]*basicModel.Apply, error) {
	return basicDao.ApplyDao.ApplyGroupList(db.RDB, userId)
}

func createContactIfNotExist(tx *gorm.DB, userID, peerID string) error {
	if !basicDao.ContactsDao.IsFriend(tx, userID, peerID) {
		return basicDao.ContactsDao.Create(tx, &basicModel.Contacts{
			ID:       utils.GenerateSfIDString(),
			UserID:   userID,
			PeerId:   peerID,
			PeerType: constant.ContactsPeerType.Friend,
		})
	}
	return nil
}
