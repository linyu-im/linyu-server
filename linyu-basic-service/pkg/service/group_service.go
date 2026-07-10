package service

import (
	"errors"
	"strings"

	basicDao "github.com/linyu-im/linyu-server/linyu-basic-service/internal/dao"
	basicModel "github.com/linyu-im/linyu-server/linyu-basic-service/pkg/model"
	basicParam "github.com/linyu-im/linyu-server/linyu-basic-service/pkg/param"
	"github.com/linyu-im/linyu-server/linyu-basic-service/pkg/result"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/constant"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/db"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/request"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/response"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/utils"
	"gorm.io/gorm"
)

var GroupService = newGroupService()

func newGroupService() *groupService {
	return &groupService{}
}

type groupService struct{}

func (s groupService) GroupCreate(userId string, param *basicParam.GroupCreateParam) (*basicModel.Group, error) {
	// 邀请的用户必须是好友
	var friendIds = []string{userId}
	for _, friendId := range param.GroupMemberList {
		if is := basicDao.ContactsDao.IsFriendBothAnd(db.RDB, userId, friendId); is {
			friendIds = append(friendIds, friendId)
		}
	}
	if len(friendIds) == 1 {
		return nil, errors.New("basic.group.no-friend-selected")
	}
	// 查询当前最大群号
	maxNumber, _ := basicDao.GroupDao.GetMaxGroupNumber(db.RDB)
	// 生成纯数字唯一群号
	number := utils.GenerateGroupNumber(maxNumber, func(number string) bool {
		user := basicDao.GroupDao.GetGroupByGroupNumber(db.RDB, number)
		return user == nil
	})
	groupId := utils.GenerateSfIDString()
	group := &basicModel.Group{
		ID:            groupId,
		Name:          buildGroupNameFromFriendIds(friendIds),
		CreatorUserID: userId,
		OwnerUserID:   userId,
		GroupNumber:   number,
		MemberNum:     len(friendIds),
	}
	// 创建群聊
	err := db.RDB.Transaction(func(tx *gorm.DB) error {
		// 群成员关系新建
		for _, id := range friendIds {
			// 判断角色
			memberRole := constant.MemberRole.Member
			if userId == id {
				memberRole = constant.MemberRole.Admin
			}
			member := &basicModel.GroupMember{
				ID:         utils.GenerateSfIDString(),
				GroupID:    group.ID,
				UserID:     id,
				MemberRole: memberRole,
			}
			if err := basicDao.GroupMemberDao.Create(tx, member); err != nil {
				return err
			}
			if err := createGroupContactIfNotExist(tx, id, group.ID); err != nil {
				return err
			}
		}
		// 新建群
		if err := basicDao.GroupDao.Create(tx, group); err != nil {
			return err
		}
		return nil
	})
	return group, err
}

func buildGroupNameFromFriendIds(friendIds []string) string {
	limit := min(3, len(friendIds))
	names := make([]string, 0, limit)
	for _, id := range friendIds[:limit] {
		user := basicDao.UserDao.GetUserById(db.RDB, id)
		if user != nil && user.Username != "" {
			names = append(names, user.Username)
		}
	}
	return strings.Join(names, "、")
}

func createGroupContactIfNotExist(tx *gorm.DB, userID, groupID string) error {
	if basicDao.ContactsDao.IsGroupContact(tx, userID, groupID) {
		return nil
	}
	return basicDao.ContactsDao.Create(tx, &basicModel.Contacts{
		ID:       utils.GenerateSfIDString(),
		UserID:   userID,
		PeerId:   groupID,
		PeerType: constant.ContactsPeerType.Group,
	})
}

func (s groupService) GroupDissolve(userId string, param *basicParam.GroupDissolveParam) error {
	// 验收是否是群主
	if !s.IsOwnerUser(param.GroupId, userId) {
		return errors.New("param.error")
	}
	err := db.RDB.Transaction(func(tx *gorm.DB) error {
		// 清空群成员
		if err := basicDao.GroupMemberDao.UnscopedDeleteMemberByGroupId(tx, param.GroupId); err != nil {
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

func (s groupService) IsOwnerUser(groupId string, userId string) bool {
	group := basicDao.GroupDao.GetGroupById(db.RDB, groupId)
	if group == nil || group.OwnerUserID != userId {
		return false
	}
	return true
}

func (s groupService) IsGroupMember(groupId string, userId string) bool {
	member := basicDao.GroupMemberDao.GetGroupMemberByGroupIdAndUserId(db.RDB, groupId, userId)
	if member == nil {
		return false
	}
	return true
}

func (s groupService) InviteMember(userId string, param *basicParam.GroupInviteMemberParam) error {
	// 验证是否是群成员
	if !s.IsGroupMember(param.GroupId, userId) {
		return errors.New("param.error")
	}
	// 过滤非好友用户和已经存在的用户
	var friendIds []string
	for _, friendId := range param.GroupMemberList {
		if is := basicDao.ContactsDao.IsContactByUserAndPeer(db.RDB, userId, friendId); !is {
			continue
		}
		if is := s.IsGroupMember(param.GroupId, friendId); is {
			continue
		}
		friendIds = append(friendIds, friendId)
	}
	err := db.RDB.Transaction(func(tx *gorm.DB) error {
		// 群成员关系新建
		for _, id := range friendIds {
			err := basicDao.GroupMemberDao.Create(tx, &basicModel.GroupMember{
				ID:         utils.GenerateSfIDString(),
				GroupID:    param.GroupId,
				UserID:     id,
				MemberRole: constant.MemberRole.Member,
			})
			if err != nil {
				return err
			}
		}
		//更新群员数量
		if err := basicDao.GroupDao.UpdateMemberNum(tx, param.GroupId); err != nil {
			return err
		}
		return nil
	})
	return err
}

func (s groupService) RemoveMember(userId string, param *basicParam.GroupRemoveMemberParam) error {
	// 验证是否是群管理
	if !s.isGroupRole(param.GroupId, userId, constant.MemberRole.Admin) {
		return errors.New("param.error")
	}
	err := db.RDB.Transaction(func(tx *gorm.DB) error {
		for _, id := range param.GroupMemberList {
			// 移除普通成员
			if s.isGroupRole(param.GroupId, id, constant.MemberRole.Member) {
				if err := basicDao.GroupMemberDao.UnscopedRemoveMember(tx, param.GroupId, id); err != nil {
					return err
				}
			}
		}
		//更新群员数量
		if err := basicDao.GroupDao.UpdateMemberNum(tx, param.GroupId); err != nil {
			return err
		}
		return nil
	})
	return err
}

func (s groupService) SetAdmin(userId string, param *basicParam.GroupSetAdminParam) error {
	if !s.IsOwnerUser(param.GroupId, userId) {
		return errors.New("param.error")
	}
	// 添加管理员
	for _, id := range param.AddAdminList {
		if !s.IsGroupMember(param.GroupId, id) {
			continue
		}
		if err := basicDao.GroupMemberDao.UpdateMemberRole(db.RDB, param.GroupId, id, constant.MemberRole.Admin); err != nil {
			return err
		}
	}
	// 移除管理员
	for _, id := range param.RemoveAdminList {
		if !s.isGroupRole(param.GroupId, id, constant.MemberRole.Admin) {
			continue
		}
		if err := basicDao.GroupMemberDao.UpdateMemberRole(db.RDB, param.GroupId, id, constant.MemberRole.Member); err != nil {
			return err
		}
	}
	return nil
}

func (s groupService) isGroupRole(groupId string, userId string, role string) bool {
	member := basicDao.GroupMemberDao.GetGroupMemberByGroupIdAndUserId(db.RDB, groupId, userId)
	if member == nil {
		return false
	}
	return member.MemberRole == role
}

func (s groupService) IsGroupAdmin(groupId string, userId string) bool {
	return s.isGroupRole(groupId, userId, constant.MemberRole.Admin)
}

func (s groupService) GetMemberUserIdsByGroupId(groupId string) []string {
	members, err := basicDao.GroupMemberDao.GetMemberUserIdsByGroupId(db.RDB, groupId)
	if err != nil {
		return []string{}
	}
	return members
}

func (s groupService) SearchByKeyword(param *basicParam.GroupSearchParam) (*response.PageResult[*basicModel.Group], error) {
	tx := basicDao.GroupDao.SearchByKeyword(db.RDB, param.Keyword)
	return response.Paginate[*basicModel.Group](tx, param.PageQuery)
}

func (s groupService) UpdateAvatar(userId string, groupId string, avatarUrl string) error {
	if !s.IsOwnerUser(groupId, userId) {
		return errors.New("param.error")
	}
	return basicDao.GroupDao.UpdateAvatar(db.RDB, groupId, avatarUrl)
}

func (s groupService) GetGroupAvatar(groupId string) interface{} {
	group := basicDao.GroupDao.GetGroupById(db.RDB, groupId)
	if group == nil {
		return ""
	}
	return group.Avatar
}

func (s groupService) GroupMemberList(userId string, param *basicParam.GroupMemberListParam) ([]*basicModel.GroupMember, error) {
	if !s.IsGroupMember(param.GroupId, userId) {
		return nil, errors.New("param.error")
	}
	tx := basicDao.GroupMemberDao.BuildGroupMemberQuery(db.RDB, param.GroupId)
	tx = tx.Order("group_user_level DESC")
	var members []*basicModel.GroupMember
	if err := tx.Find(&members).Error; err != nil {
		return nil, err
	}
	return members, nil
}

func (s groupService) UpdateInfo(userId string, param *basicParam.GroupUpdateInfoParam) error {
	if !s.IsOwnerUser(param.GroupId, userId) {
		return errors.New("param.error")
	}
	fields := map[string]interface{}{}
	if param.Name != "" {
		fields["name"] = param.Name
	}
	fields["describe"] = param.Describe
	fields["tag"] = param.Tag
	return basicDao.GroupDao.UpdateInfo(db.RDB, param.GroupId, fields)
}

func (s groupService) GroupInfo(userId string, groupId string) *result.GroupInfoResult {
	group := basicDao.GroupDao.GroupInfoById(db.RDB, userId, groupId)
	//获取top6群成员
	tx := basicDao.GroupMemberDao.BuildGroupMemberQuery(db.RDB, groupId)
	pages, err := response.Paginate[*basicModel.GroupMember](tx, request.PageQuery{
		Page:      1,
		PageSize:  6,
		SortBy:    "group_user_level",
		SortOrder: "desc",
	})
	if err != nil {
		return nil
	}
	return &result.GroupInfoResult{
		Info: group,
		Tops: pages.Records,
	}
}
