package dao

import (
	driveModel "github.com/linyu-im/linyu-server/linyu-cloud-drive/pkg/model"
	"gorm.io/gorm"
)

var SpaceDao = newSpaceDao()

func newSpaceDao() *spaceDao {
	return &spaceDao{}
}

type spaceDao struct{}

func (d *spaceDao) GetByTypeAndTargetID(db *gorm.DB, spaceType string, targetId string) *driveModel.Space {
	result := &driveModel.Space{}
	if err := db.Where("space_type = ? AND target_id = ?", spaceType, targetId).First(result).Error; err != nil {
		return nil
	}
	return result
}

func (d *spaceDao) GetById(db *gorm.DB, id string) *driveModel.Space {
	result := &driveModel.Space{}
	if err := db.First(result, "id = ?", id).Error; err != nil {
		return nil
	}
	return result
}

func (d *spaceDao) Create(db *gorm.DB, space *driveModel.Space) error {
	return db.Create(space).Error
}
