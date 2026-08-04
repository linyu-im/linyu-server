package dao

import (
	driveModel "github.com/linyu-im/linyu-server/linyu-cloud-drive/pkg/model"
	driveResult "github.com/linyu-im/linyu-server/linyu-cloud-drive/pkg/result"
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

func (d *spaceFileDao) ListDirsBySpaceId(db *gorm.DB, spaceId string) ([]*driveModel.SpaceFile, error) {
	var list []*driveModel.SpaceFile
	if err := db.Where("space_id = ? AND is_dir = ?", spaceId, true).
		Order("level ASC, file_name ASC").
		Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

// ListFilesUnderPath 查询指定目录 path 下全部文件（含子目录，不含目录）
// pathPrefix 为空表示空间根下全部文件
func (d *spaceFileDao) ListFilesUnderPath(db *gorm.DB, spaceId, pathPrefix string) ([]*driveModel.SpaceFile, error) {
	var list []*driveModel.SpaceFile
	query := db.Where("space_id = ? AND is_dir = ?", spaceId, false)
	if pathPrefix != "" {
		query = query.Where("path LIKE ?", pathPrefix+"/%")
	}
	if err := query.Order("updated_at DESC").Find(&list).Error; err != nil {
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

func (d *spaceFileDao) ListByIdsUnscoped(db *gorm.DB, ids []string) ([]*driveModel.SpaceFile, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	var list []*driveModel.SpaceFile
	if err := db.Unscoped().Where("id IN ?", ids).Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (d *spaceFileDao) SumFileSizeSelfAndDescendants(db *gorm.DB, spaceId, id, path string) (int64, int64, error) {
	type row struct {
		TotalSize int64
		FileCount int64
	}
	var r row
	err := db.Model(&driveModel.SpaceFile{}).
		Select("COALESCE(SUM(file_size), 0) AS total_size, COUNT(*) AS file_count").
		Where("space_id = ? AND is_dir = ? AND (id = ? OR path LIKE ?)", spaceId, false, id, path+"/%").
		Scan(&r).Error
	if err != nil {
		return 0, 0, err
	}
	return r.TotalSize, r.FileCount, nil
}

// CountDescendants 统计目录下子孙文件数与文件夹数
func (d *spaceFileDao) CountDescendants(db *gorm.DB, spaceId, path string) (fileCount, folderCount int64, err error) {
	type row struct {
		FileCount   int64
		FolderCount int64
	}
	var r row
	err = db.Model(&driveModel.SpaceFile{}).
		Select(`
			COALESCE(SUM(CASE WHEN is_dir = 0 THEN 1 ELSE 0 END), 0) AS file_count,
			COALESCE(SUM(CASE WHEN is_dir = 1 THEN 1 ELSE 0 END), 0) AS folder_count
		`).
		Where("space_id = ? AND path LIKE ?", spaceId, path+"/%").
		Scan(&r).Error
	if err != nil {
		return 0, 0, err
	}
	return r.FileCount, r.FolderCount, nil
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

// ListSelfAndDescendants 查询自身及全部子孙
func (d *spaceFileDao) ListSelfAndDescendants(db *gorm.DB, spaceId string, id string, path string) ([]*driveModel.SpaceFile, error) {
	var list []*driveModel.SpaceFile
	if err := db.Where("space_id = ? AND (id = ? OR path LIKE ?)", spaceId, id, path+"/%").
		Order("level ASC").
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

func (d *spaceFileDao) UpdateMove(db *gorm.DB, id, parentId, newPath string, newLevel int) error {
	return db.Model(&driveModel.SpaceFile{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"parent_id": parentId,
			"path":      newPath,
			"level":     newLevel,
		}).Error
}

func (d *spaceFileDao) UpdateName(db *gorm.DB, id, fileName string, fileType, fileCategory string, updateTypeMeta bool) error {
	updates := map[string]interface{}{
		"file_name": fileName,
	}
	if updateTypeMeta {
		updates["file_type"] = fileType
		updates["file_category"] = fileCategory
	}
	return db.Model(&driveModel.SpaceFile{}).
		Where("id = ?", id).
		Updates(updates).Error
}

// UpdateDescendantsPathAndLevel 移动目录后，批量更新所有子孙节点 path/level
func (d *spaceFileDao) UpdateDescendantsPathAndLevel(db *gorm.DB, spaceId, oldPath, newPath string, levelDelta int) error {
	return db.Model(&driveModel.SpaceFile{}).
		Where("space_id = ? AND path LIKE ?", spaceId, oldPath+"/%").
		Updates(map[string]interface{}{
			"path":  gorm.Expr("CONCAT(?, SUBSTRING(path, ?))", newPath, len(oldPath)+1),
			"level": gorm.Expr("level + ?", levelDelta),
		}).Error
}

func (d *spaceFileDao) StatByCategory(db *gorm.DB, spaceId string) ([]*driveResult.SpaceFileCategoryStat, error) {
	var list []*driveResult.SpaceFileCategoryStat
	if err := db.Model(&driveModel.SpaceFile{}).
		Select("file_category AS file_category, COUNT(*) AS file_count, COALESCE(SUM(file_size), 0) AS total_size").
		Where("space_id = ? AND is_dir = ?", spaceId, false).
		Group("file_category").
		Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}
