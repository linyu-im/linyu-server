package service

import (
	"errors"
	"strings"
	"time"

	driveDao "github.com/linyu-im/linyu-server/linyu-cloud-drive/internal/dao"
	driveModel "github.com/linyu-im/linyu-server/linyu-cloud-drive/pkg/model"
	driveResult "github.com/linyu-im/linyu-server/linyu-cloud-drive/pkg/result"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/constant"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/db"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/localtime"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/utils"
	"gorm.io/gorm"
)

var SpaceService = newSpaceService()

func newSpaceService() *spaceService {
	return &spaceService{}
}

type spaceService struct{}

// CheckUserSpaceQuota 校验用户空间容量是否足够容纳本次上传文件
func (s *spaceService) CheckUserSpaceQuota(userId string, fileSize int64) error {
	space, err := s.GetOrCreateUserSpace(userId)
	if err != nil {
		return err
	}
	return s.CheckSpaceQuota(space, fileSize)
}

// CheckSpaceQuota 校验空间容量是否足够
func (s *spaceService) CheckSpaceQuota(space *driveModel.Space, fileSize int64) error {
	// QuotaBytes 为 0 表示不限制
	if space.QuotaBytes > 0 && space.UsedBytes+fileSize > space.QuotaBytes {
		return errors.New("cloud-drive.space.quota-exceeded")
	}
	return nil
}

// CalcSpaceFilesUsedBytes 计算文件/文件夹占用容量（文件夹通过 path like 汇总全部子文件）
func (s *spaceService) CalcSpaceFilesUsedBytes(tx *gorm.DB, spaceId string, files []*driveModel.SpaceFile) (int64, int64, error) {
	topFiles := filterTopLevelSpaceFiles(files)
	var totalSize, fileCount int64
	for _, file := range topFiles {
		if file.IsDir {
			size, count, err := driveDao.SpaceFileDao.SumFileSizeSelfAndDescendants(tx, spaceId, file.ID, file.Path)
			if err != nil {
				return 0, 0, err
			}
			totalSize += size
			fileCount += count
			continue
		}
		totalSize += file.FileSize
		fileCount++
	}
	return totalSize, fileCount, nil
}

// filterTopLevelSpaceFiles 过滤掉已被其他选中文件夹覆盖的项，避免重复统计
func filterTopLevelSpaceFiles(files []*driveModel.SpaceFile) []*driveModel.SpaceFile {
	result := make([]*driveModel.SpaceFile, 0, len(files))
	for _, file := range files {
		covered := false
		for _, other := range files {
			if other.ID == file.ID || !other.IsDir {
				continue
			}
			if strings.HasPrefix(file.Path, other.Path+"/") {
				covered = true
				break
			}
		}
		if !covered {
			result = append(result, file)
		}
	}
	return result
}

func (s *spaceService) GetOrCreateUserSpace(userId string) (*driveModel.Space, error) {
	if space := driveDao.SpaceDao.GetByTypeAndTargetID(db.RDB, constant.SpaceType.User, userId); space != nil {
		return space, nil
	}

	spaceID := utils.GenerateSfIDString()
	space := &driveModel.Space{
		ID:         spaceID,
		SpaceType:  constant.SpaceType.User,
		TargetID:   userId,
		OwnerID:    userId,
		SpaceName:  "我的云盘",
		QuotaBytes: constant.DefaultUserSpaceQuotaBytes,
		Status:     constant.SpaceStatus.Active,
	}

	err := db.RDB.Transaction(func(tx *gorm.DB) error {
		if err := driveDao.SpaceDao.Create(tx, space); err != nil {
			return err
		}
		return driveDao.SpaceMemberDao.Create(tx, &driveModel.SpaceMember{
			ID:         utils.GenerateSfIDString(),
			SpaceID:    spaceID,
			UserID:     userId,
			MemberRole: constant.SpaceMemberRole.Owner,
			Status:     constant.SpaceMemberStatus.Active,
		})
	})
	if err != nil {
		if exist := driveDao.SpaceDao.GetByTypeAndTargetID(db.RDB, constant.SpaceType.User, userId); exist != nil {
			return exist, nil
		}
		return nil, err
	}
	return space, nil
}

func (s *spaceService) ListUserSpaceFiles(userId string, parentId string) ([]*driveModel.SpaceFile, error) {
	space, err := s.GetOrCreateUserSpace(userId)
	if err != nil {
		return nil, err
	}
	if _, err := SpaceFileService.VerifyDirPermission(parentId, space.ID, userId); err != nil {
		return nil, err
	}
	return driveDao.SpaceFileDao.ListBySpaceIdAndParentId(db.RDB, space.ID, parentId)
}

func (s *spaceService) CreateUserSpaceDir(userId string, parentId string, dirName string) (*driveModel.SpaceFile, error) {
	space, err := s.GetOrCreateUserSpace(userId)
	if err != nil {
		return nil, err
	}
	parentFile, err := SpaceFileService.VerifyDirPermission(parentId, space.ID, userId)
	if err != nil {
		return nil, err
	}

	path := ""
	level := 0
	resolvedParentId := "root"
	if parentFile != nil {
		path = parentFile.Path
		level = parentFile.Level
		resolvedParentId = parentFile.ID
	}
	id := utils.GenerateSfIDString()
	dir := &driveModel.SpaceFile{
		ID:                  id,
		SpaceID:             space.ID,
		UserID:              userId,
		PhysicalID:          "",
		PhysicalStoragePath: "",
		ParentID:            resolvedParentId,
		Path:                path + "/" + id,
		Level:               level + 1,
		FileName:            dirName,
		IsDir:               true,
		FileType:            "",
		FileSize:            0,
	}
	if err := driveDao.SpaceFileDao.Create(db.RDB, dir); err != nil {
		return nil, err
	}
	return dir, nil
}

func (s *spaceService) DeleteUserSpaceFiles(userId string, spaceFileIds []string) error {
	space, err := s.GetOrCreateUserSpace(userId)
	if err != nil {
		return err
	}

	idSet := make(map[string]struct{}, len(spaceFileIds))
	uniqueIds := make([]string, 0, len(spaceFileIds))
	for _, id := range spaceFileIds {
		if id == "" {
			return errors.New("param.error")
		}
		if _, ok := idSet[id]; ok {
			continue
		}
		idSet[id] = struct{}{}
		uniqueIds = append(uniqueIds, id)
	}

	files, err := driveDao.SpaceFileDao.ListByIds(db.RDB, uniqueIds)
	if err != nil {
		return err
	}
	if len(files) != len(uniqueIds) {
		return errors.New("param.error")
	}
	for _, file := range files {
		if file.SpaceID != space.ID {
			return errors.New("param.error")
		}
	}

	totalSize, fileCount, err := s.CalcSpaceFilesUsedBytes(db.RDB, space.ID, files)
	if err != nil {
		return err
	}

	expireAt := localtime.LocalTime(time.Now().In(localtime.Location).
		AddDate(0, 0, constant.DefaultSpaceRecycleExpireDays))
	ids := make([]string, 0, len(files))
	recycles := make([]*driveModel.SpaceRecycle, 0, len(files))
	for _, file := range files {
		ids = append(ids, file.ID)
		recycles = append(recycles, &driveModel.SpaceRecycle{
			ID:          utils.GenerateSfIDString(),
			UserID:      userId,
			SpaceID:     space.ID,
			SpaceFileID: file.ID,
			DeletedBy:   userId,
			ExpireAt:    &expireAt,
		})
	}

	return db.RDB.Transaction(func(tx *gorm.DB) error {
		if err := driveDao.SpaceFileDao.DeleteByIds(tx, ids); err != nil {
			return err
		}
		if err := driveDao.SpaceRecycleDao.CreateBatch(tx, recycles); err != nil {
			return err
		}
		if totalSize == 0 && fileCount == 0 {
			return nil
		}
		return driveDao.SpaceDao.DecUsedBytesById(tx, space.ID, totalSize, fileCount)
	})
}

func (s *spaceService) ListUserSpaceCategoryStats(userId string) ([]*driveResult.SpaceFileCategoryStat, error) {
	space, err := s.GetOrCreateUserSpace(userId)
	if err != nil {
		return nil, err
	}
	stats, err := driveDao.SpaceFileDao.StatByCategory(db.RDB, space.ID)
	if err != nil {
		return nil, err
	}

	statMap := make(map[string]*driveResult.SpaceFileCategoryStat, len(stats))
	for _, stat := range stats {
		statMap[stat.FileCategory] = stat
	}

	categories := []string{
		constant.FileCategory.Image,
		constant.FileCategory.Video,
		constant.FileCategory.Document,
		constant.FileCategory.Audio,
		constant.FileCategory.Archive,
		constant.FileCategory.Other,
	}
	result := make([]*driveResult.SpaceFileCategoryStat, 0, len(categories))
	for _, category := range categories {
		if stat, ok := statMap[category]; ok {
			result = append(result, stat)
			continue
		}
		result = append(result, &driveResult.SpaceFileCategoryStat{
			FileCategory: category,
			FileCount:    0,
			TotalSize:    0,
		})
	}
	return result, nil
}
