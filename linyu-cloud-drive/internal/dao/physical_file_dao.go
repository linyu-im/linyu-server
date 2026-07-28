package dao

import (
	driveModel "github.com/linyu-im/linyu-server/linyu-cloud-drive/pkg/model"
	"gorm.io/gorm"
)

var PhysicalFileDao = newPhysicalFileDao()

func newPhysicalFileDao() *physicalFileDao {
	return &physicalFileDao{}
}

type physicalFileDao struct{}

func (d *physicalFileDao) GetByHash(db *gorm.DB, hash string) *driveModel.PhysicalFile {
	result := &driveModel.PhysicalFile{}
	if err := db.First(result, "file_hash = ?", hash).Error; err != nil {
		return nil
	}
	return result
}

func (d *physicalFileDao) FileRefIncById(db *gorm.DB, id string) error {
	err := db.Model(&driveModel.PhysicalFile{}).
		Where("id = ?", id).
		Update("ref_count", gorm.Expr("ref_count + ?", 1)).Error
	return err
}

func (d *physicalFileDao) FileRefDecById(db *gorm.DB, id string) error {
	if id == "" {
		return nil
	}
	return db.Model(&driveModel.PhysicalFile{}).
		Where("id = ? AND ref_count > 0", id).
		Update("ref_count", gorm.Expr("ref_count - ?", 1)).Error
}

func (d *physicalFileDao) Create(db *gorm.DB, file *driveModel.PhysicalFile) error {
	if err := db.Create(file).Error; err != nil {
		return err
	}
	return nil
}
