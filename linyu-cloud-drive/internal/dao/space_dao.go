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

func (d *spaceDao) IncUsedBytesById(db *gorm.DB, id string, size int64, fileCount int64) error {
	return db.Model(&driveModel.Space{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"used_bytes": gorm.Expr("used_bytes + ?", size),
			"file_count": gorm.Expr("file_count + ?", fileCount),
		}).Error
}

func (d *spaceDao) DecUsedBytesById(db *gorm.DB, id string, size int64, fileCount int64) error {
	return db.Model(&driveModel.Space{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"used_bytes": gorm.Expr("GREATEST(used_bytes - ?, 0)", size),
			"file_count": gorm.Expr("GREATEST(file_count - ?, 0)", fileCount),
		}).Error
}
