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

func (d *groupMemberDao) DeleteMemberByGroupId(db *gorm.DB, groupId string) error {
	err := db.Where("group_id = ?", groupId).Delete(&basicModel.GroupMember{}).Error
	return err
}
