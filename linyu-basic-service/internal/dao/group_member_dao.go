package dao

import (
	basicModel "github.com/linyu-im/linyu-server/linyu-basic-service/pkg/model"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/constant"
	"gorm.io/gorm"
)

var GroupMemberDao = newGroupMemberDao()

func newGroupMemberDao() *groupMemberDao {
	return &groupMemberDao{}
}

type groupMemberDao struct{}

func (d *groupMemberDao) Create(db *gorm.DB, groupMember *basicModel.GroupMember) error {
	if err := db.Create(groupMember).Error; err != nil {
		return err
	}
	return nil
}

func (d *groupMemberDao) UnscopedDeleteMemberByGroupId(db *gorm.DB, groupId string) error {
	err := db.Unscoped().Where("group_id = ?", groupId).Delete(&basicModel.GroupMember{}).Error
	return err
}

func (d *groupMemberDao) GetGroupMemberByGroupIdAndUserId(db *gorm.DB, groupId string, userId string) *basicModel.GroupMember {
	result := &basicModel.GroupMember{}
	if err := db.First(result, "group_id = ? AND user_id = ?", groupId, userId).Error; err != nil {
		return nil
	}
	return result
}

func (d *groupMemberDao) UnscopedRemoveMember(db *gorm.DB, groupId string, userId string) error {
	err := db.Unscoped().Where("group_id = ? AND user_id = ?", groupId, userId).Delete(&basicModel.GroupMember{}).Error
	return err
}

func (d *groupMemberDao) GetMembersByGroupId(db *gorm.DB, groupId string) ([]*basicModel.GroupMember, error) {
	var members []*basicModel.GroupMember
	if err := db.Where("group_id = ?", groupId).Find(&members).Error; err != nil {
		return nil, err
	}
	return members, nil
}

func (d *groupMemberDao) GetMemberUserIdsByGroupId(db *gorm.DB, groupId string) ([]string, error) {
	var userIds []string
	err := db.Model(&basicModel.GroupMember{}).
		Where("group_id = ?", groupId).
		Pluck("user_id", &userIds).
		Error
	if err != nil {
		return nil, err
	}
	return userIds, nil
}

func (d *groupMemberDao) GetAdminUserIdsByGroupId(db *gorm.DB, groupId string) []string {
	var userIds []string
	db.Model(&basicModel.GroupMember{}).
		Where("group_id = ? AND member_role = ?", groupId, constant.MemberRole.Admin).
		Pluck("user_id", &userIds)
	return userIds
}

func (d *groupMemberDao) UpdateMemberRole(db *gorm.DB, groupId string, userId string, role string) error {
	return db.Model(&basicModel.GroupMember{}).
		Where("group_id = ? AND user_id = ?", groupId, userId).
		Update("member_role", role).Error
}

func (d *groupMemberDao) UpdateGroupNickName(db *gorm.DB, groupId string, userId string, nickName string) error {
	return db.Model(&basicModel.GroupMember{}).
		Where("group_id = ? AND user_id = ?", groupId, userId).
		Update("group_nick_name", nickName).Error
}

func (d *groupMemberDao) BuildGroupMemberQuery(db *gorm.DB, id string) *gorm.DB {
	return db.Table("t_group_member AS gm").
		Select("gm.*, u.username, e.emotion_name AS emotion_name, e.url AS emotion_url").
		Joins("LEFT JOIN t_user AS u ON gm.user_id = u.id").
		Joins("LEFT JOIN t_emotion AS e ON u.emotion_id = e.id").
		Where("group_id = ?", id)
}
