package dao

import (
	basicModel "github.com/linyu-im/linyu-server/linyu-basic-service/pkg/model"
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
