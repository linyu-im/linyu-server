package service

import (
	"errors"

	driveDao "github.com/linyu-im/linyu-server/linyu-cloud-drive/internal/dao"
	driveModel "github.com/linyu-im/linyu-server/linyu-cloud-drive/pkg/model"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/db"
	"gorm.io/gorm"
)

var SpaceRecycleService = newSpaceRecycleService()

func newSpaceRecycleService() *spaceRecycleService {
	return &spaceRecycleService{}
}

type spaceRecycleService struct{}

func (s *spaceRecycleService) ListUserRecycle(userId string) ([]*driveModel.SpaceRecycle, error) {
	space, err := SpaceService.GetOrCreateUserSpace(userId)
	if err != nil {
		return nil, err
	}
	return driveDao.SpaceRecycleDao.ListByUserIdAndSpaceId(db.RDB, userId, space.ID)
}

func (s *spaceRecycleService) RestoreUserRecycle(userId string, spaceRecycleIds []string) error {
	space, err := SpaceService.GetOrCreateUserSpace(userId)
	if err != nil {
		return err
	}
	recycles, err := driveDao.SpaceRecycleDao.ListByIdsAndUserIdAndSpaceId(db.RDB, spaceRecycleIds, userId, space.ID)
	if err != nil {
		return err
	}
	if len(recycles) == 0 {
		return errors.New("param.error")
	}

	spaceFileIds := make([]string, len(recycles))
	recycleIds := make([]string, len(recycles))
	for i, r := range recycles {
		spaceFileIds[i] = r.SpaceFileID
		recycleIds[i] = r.ID
	}

	files, err := driveDao.SpaceFileDao.ListByIdsUnscoped(db.RDB, spaceFileIds)
	if err != nil {
		return err
	}
	if len(files) != len(spaceFileIds) {
		return errors.New("param.error")
	}
	for _, file := range files {
		if file.SpaceID != space.ID {
			return errors.New("param.error")
		}
	}

	// 文件夹通过 path like 汇总子文件占用；还原前校验容量是否够用
	totalSize, fileCount, err := SpaceService.CalcSpaceFilesUsedBytes(db.RDB.Unscoped(), space.ID, files)
	if err != nil {
		return err
	}
	if err := SpaceService.CheckSpaceQuota(space, totalSize); err != nil {
		return err
	}

	return db.RDB.Transaction(func(tx *gorm.DB) error {
		if err := driveDao.SpaceFileDao.ClearDeletedAtByIds(tx, spaceFileIds); err != nil {
			return err
		}
		if err := driveDao.SpaceRecycleDao.UnscopedDeleteByIds(tx, recycleIds); err != nil {
			return err
		}
		if totalSize == 0 && fileCount == 0 {
			return nil
		}
		return driveDao.SpaceDao.IncUsedBytesById(tx, space.ID, totalSize, fileCount)
	})
}

func (s *spaceRecycleService) PermanentlyDeleteUserRecycle(userId string, spaceRecycleIds []string) error {
	space, err := SpaceService.GetOrCreateUserSpace(userId)
	if err != nil {
		return err
	}
	recycles, err := driveDao.SpaceRecycleDao.ListByIdsAndUserIdAndSpaceIdIgnoreExpire(db.RDB, spaceRecycleIds, userId, space.ID)
	if err != nil {
		return err
	}
	if len(recycles) == 0 {
		return errors.New("param.error")
	}

	fileIDSet := make(map[string]struct{})
	var deleteFiles []*driveModel.SpaceFile
	recycleIds := make([]string, 0, len(recycles))

	for _, recycle := range recycles {
		recycleIds = append(recycleIds, recycle.ID)
		file := driveDao.SpaceFileDao.GetByIdUnscoped(db.RDB, recycle.SpaceFileID)
		if file == nil || file.SpaceID != space.ID {
			return errors.New("param.error")
		}

		var targets []*driveModel.SpaceFile
		if file.IsDir {
			targets, err = driveDao.SpaceFileDao.ListSelfAndDescendantsUnscoped(db.RDB, space.ID, file.ID, file.Path)
			if err != nil {
				return err
			}
		} else {
			targets = []*driveModel.SpaceFile{file}
		}

		for _, target := range targets {
			if _, ok := fileIDSet[target.ID]; ok {
				continue
			}
			fileIDSet[target.ID] = struct{}{}
			deleteFiles = append(deleteFiles, target)
		}
	}

	spaceFileIds := make([]string, 0, len(deleteFiles))
	physicalIds := make([]string, 0)
	for _, file := range deleteFiles {
		spaceFileIds = append(spaceFileIds, file.ID)
		if !file.IsDir && file.PhysicalID != "" {
			physicalIds = append(physicalIds, file.PhysicalID)
		}
	}

	return db.RDB.Transaction(func(tx *gorm.DB) error {
		for _, physicalId := range physicalIds {
			if err := driveDao.PhysicalFileDao.FileRefDecById(tx, physicalId); err != nil {
				return err
			}
		}
		if err := driveDao.SpaceFileDao.UnscopedDeleteByIds(tx, spaceFileIds); err != nil {
			return err
		}
		return driveDao.SpaceRecycleDao.UnscopedDeleteByIds(tx, recycleIds)
	})
}

func (s *spaceRecycleService) ClearUserRecycle(userId string) error {
	space, err := SpaceService.GetOrCreateUserSpace(userId)
	if err != nil {
		return err
	}
	ids, err := driveDao.SpaceRecycleDao.ListIdsByUserIdAndSpaceId(db.RDB, userId, space.ID)
	if err != nil {
		return err
	}
	if len(ids) == 0 {
		return nil
	}
	return s.PermanentlyDeleteUserRecycle(userId, ids)
}
