package dao

import (
	driveModel "github.com/linyu-im/linyu-server/linyu-cloud-drive/pkg/model"
	"gorm.io/gorm"
)

var SpaceMemberDao = newSpaceMemberDao()

func newSpaceMemberDao() *spaceMemberDao {
	return &spaceMemberDao{}
}

type spaceMemberDao struct{}

func (d *spaceMemberDao) Create(db *gorm.DB, member *driveModel.SpaceMember) error {
	return db.Create(member).Error
}
