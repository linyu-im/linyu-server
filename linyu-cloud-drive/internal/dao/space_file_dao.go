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

func (d *spaceFileDao) GetByIdUnscoped(db *gorm.DB, id string) *driveModel.SpaceFile {
	result := &driveModel.SpaceFile{}
	if err := db.Unscoped().First(result, "id = ?", id).Error; err != nil {
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

func (d *spaceFileDao) ListBySpaceIdAndParentId(db *gorm.DB, spaceId string, parentId string) ([]*driveModel.SpaceFile, error) {
	var list []*driveModel.SpaceFile
	if err := db.Where("space_id = ? AND parent_id = ?", spaceId, parentId).
		Order("updated_at DESC").
		Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (d *spaceFileDao) DeleteById(db *gorm.DB, id string) error {
	return db.Where("id = ?", id).Delete(&driveModel.SpaceFile{}).Error
}

func (d *spaceFileDao) DeleteByIds(db *gorm.DB, ids []string) error {
	if len(ids) == 0 {
		return nil
	}
	return db.Where("id IN ?", ids).Delete(&driveModel.SpaceFile{}).Error
}

func (d *spaceFileDao) ListByIds(db *gorm.DB, ids []string) ([]*driveModel.SpaceFile, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	var list []*driveModel.SpaceFile
	if err := db.Where("id IN ?", ids).Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (d *spaceFileDao) ClearDeletedAtByIds(db *gorm.DB, ids []string) error {
	if len(ids) == 0 {
		return nil
	}
	return db.Unscoped().Model(&driveModel.SpaceFile{}).
		Where("id IN ?", ids).
		Updates(map[string]interface{}{"deleted_at": nil}).Error
}

func (d *spaceFileDao) ListSelfAndDescendantsUnscoped(db *gorm.DB, spaceId string, id string, path string) ([]*driveModel.SpaceFile, error) {
	var list []*driveModel.SpaceFile
	if err := db.Unscoped().
		Where("space_id = ? AND (id = ? OR path LIKE ?)", spaceId, id, path+"/%").
		Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (d *spaceFileDao) UnscopedDeleteByIds(db *gorm.DB, ids []string) error {
	if len(ids) == 0 {
		return nil
	}
	return db.Unscoped().Where("id IN ?", ids).Delete(&driveModel.SpaceFile{}).Error
}
