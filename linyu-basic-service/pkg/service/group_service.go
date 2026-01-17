package service

import (
	"errors"
	basicDao "github.com/linyu-im/linyu-server/linyu-basic-service/internal/dao"
	basicModel "github.com/linyu-im/linyu-server/linyu-basic-service/pkg/model"
	basicParam "github.com/linyu-im/linyu-server/linyu-basic-service/pkg/param"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/db"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/utils"
	"gorm.io/gorm"
)

var GroupService = newGroupService()

func newGroupService() *groupService {
	return &groupService{}
}

type groupService struct{}

func (s groupService) GroupCreate(userId string, param *basicParam.GroupCreateParam) error {
	// 邀请的用户必须是好友
	var friendIds = []string{userId}
	for _, friendId := range param.GroupMemberList {
		if is := basicDao.ContactsDao.IsContactByUserAndPeer(db.MysqlDB, userId, friendId); is {
			friendIds = append(friendIds, friendId)
		}
	}
	// 生成唯一群号
	number := utils.GenerateOnlyNumber("LINYU-", func(number string) bool {
		user := basicDao.GroupDao.GetGroupByGroupNumber(db.MysqlDB, number)
		return user == nil
	})
	group := &basicModel.Group{
		ID:            utils.GenerateSfIDString(),
		Name:          param.GroupName,
		CreatorUserID: userId,
		OwnerUserID:   userId,
		GroupNumber:   number,
		MemberNum:     len(friendIds),
	}
	// 创建群聊
	err := db.MysqlDB.Transaction(func(tx *gorm.DB) error {
		// 群成员关系新建
		for _, id := range friendIds {
			err := basicDao.GroupMemberDao.Create(tx, &basicModel.GroupMember{
				ID:      utils.GenerateSfIDString(),
				GroupID: group.ID,
				UserID:  id,
			})
			if err != nil {
				return err
			}
		}
		// 新建群
		if err := basicDao.GroupDao.Create(tx, group); err != nil {
			return err
		}
		return nil
	})
	return err
}

func (s groupService) GroupDissolve(userId string, param *basicParam.GroupDissolveParam) error {
	// 验收是否是群主
	if !s.isOwnerUser(param.GroupId, userId) {
		return errors.New("param.error")
	}
	err := db.MysqlDB.Transaction(func(tx *gorm.DB) error {
		// 清空群成员
		if err := basicDao.GroupMemberDao.DeleteMemberByGroupId(tx, param.GroupId); err != nil {
			return err
		}
		// 删除群聊
		if err := basicDao.GroupDao.DeleteById(tx, param.GroupId); err != nil {
			return err
		}
		return nil
	})
	return err
}

func (s groupService) isOwnerUser(groupId string, userId string) bool {
	group := basicDao.GroupDao.GetGroupById(db.MysqlDB, groupId)
	if group == nil || group.OwnerUserID != userId {
		return false
	}
	return true
}
