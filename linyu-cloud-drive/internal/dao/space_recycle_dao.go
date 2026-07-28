package dao

import (
	driveModel "github.com/linyu-im/linyu-server/linyu-cloud-drive/pkg/model"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/localtime"
	"gorm.io/gorm"
)

var SpaceRecycleDao = newSpaceRecycleDao()

func newSpaceRecycleDao() *spaceRecycleDao {
	return &spaceRecycleDao{}
}

type spaceRecycleDao struct{}

func (d *spaceRecycleDao) Create(db *gorm.DB, recycle *driveModel.SpaceRecycle) error {
	return db.Create(recycle).Error
}

func (d *spaceRecycleDao) CreateBatch(db *gorm.DB, list []*driveModel.SpaceRecycle) error {
	if len(list) == 0 {
		return nil
	}
	return db.Create(&list).Error
}

func (d *spaceRecycleDao) ListByUserIdAndSpaceId(db *gorm.DB, userId string, spaceId string) ([]*driveModel.SpaceRecycle, error) {
	var list []*driveModel.SpaceRecycle
	now := localtime.Now().ToTime().Format("2006-01-02 15:04:05")
	if err := db.Table("t_space_recycle AS r").
		Select(`r.*,
			f.file_name AS file_name,
			f.is_dir AS is_dir,
			f.file_type AS file_type,
			f.file_size AS file_size`).
		Joins("LEFT JOIN t_space_file f ON f.id = r.space_file_id").
		Where("r.user_id = ? AND r.space_id = ? AND r.deleted_at IS NULL AND r.expire_at IS NOT NULL AND r.expire_at >= ?", userId, spaceId, now).
		Order("r.created_at DESC").
		Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (d *spaceRecycleDao) ListByIdsAndUserIdAndSpaceId(db *gorm.DB, ids []string, userId string, spaceId string) ([]*driveModel.SpaceRecycle, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	var list []*driveModel.SpaceRecycle
	now := localtime.Now().ToTime().Format("2006-01-02 15:04:05")
	if err := db.Where("id IN ? AND user_id = ? AND space_id = ? AND expire_at IS NOT NULL AND expire_at >= ?",
		ids, userId, spaceId, now).
		Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (d *spaceRecycleDao) ListByIdsAndUserIdAndSpaceIdIgnoreExpire(db *gorm.DB, ids []string, userId string, spaceId string) ([]*driveModel.SpaceRecycle, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	var list []*driveModel.SpaceRecycle
	if err := db.Where("id IN ? AND user_id = ? AND space_id = ?", ids, userId, spaceId).
		Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (d *spaceRecycleDao) ListIdsByUserIdAndSpaceId(db *gorm.DB, userId string, spaceId string) ([]string, error) {
	var ids []string
	if err := db.Model(&driveModel.SpaceRecycle{}).
		Where("user_id = ? AND space_id = ?", userId, spaceId).
		Pluck("id", &ids).Error; err != nil {
		return nil, err
	}
	return ids, nil
}

func (d *spaceRecycleDao) UnscopedDeleteByIds(db *gorm.DB, ids []string) error {
	if len(ids) == 0 {
		return nil
	}
	return db.Unscoped().Where("id IN ?", ids).Delete(&driveModel.SpaceRecycle{}).Error
}
