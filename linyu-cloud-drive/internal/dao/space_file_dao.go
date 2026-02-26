package dao

import (
	driveModel "github.com/linyu-im/linyu-server/linyu-cloud-drive/pkg/model"
	"gorm.io/gorm"
)

var SpaceFileDao = newSpaceFileDao()

func newSpaceFileDao() *spaceFileDao {
	return &spaceFileDao{}
}

type spaceFileDao struct{}

func (d *spaceFileDao) GetById(db *gorm.DB, id string) *driveModel.SpaceFile {
	result := &driveModel.SpaceFile{}
	if err := db.First(result, "id = ?", id).Error; err != nil {
		return nil
	}
	return result
}

func (d *spaceFileDao) Create(db *gorm.DB, file *driveModel.SpaceFile) error {
	if err := db.Create(file).Error; err != nil {
		return err
	}
	return nil
}
