package service

import (
	"errors"
	"path/filepath"
	"strings"

	basicService "github.com/linyu-im/linyu-server/linyu-basic-service/pkg/service"
	driveDao "github.com/linyu-im/linyu-server/linyu-cloud-drive/internal/dao"
	driveModel "github.com/linyu-im/linyu-server/linyu-cloud-drive/pkg/model"
	driveParam "github.com/linyu-im/linyu-server/linyu-cloud-drive/pkg/param"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/constant"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/db"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/utils"
	"gorm.io/gorm"
)

var SpaceFileService = nweSpaceFileService()

func nweSpaceFileService() *spaceFileService {
	return &spaceFileService{}
}

type spaceFileService struct{}

func (s spaceFileService) CreateUserFileFromPhysicalFile(userId string,
	param *driveParam.UploadFileInfoParam,
	physicalFile *driveModel.PhysicalFile) error {

	space, err := SpaceService.GetOrCreateUserSpace(userId)
	if err != nil {
		return err
	}
	parentFile, err := s.VerifyDirPermission(param.ParentID, space.ID, userId)
	if err != nil {
		return err
	}
	path := ""
	level := 0
	parentId := "root"
	if parentFile != nil {
		path = parentFile.Path
		parentId = parentFile.ID
		level = parentFile.Level
	}

	err = db.RDB.Transaction(func(tx *gorm.DB) error {
		id := utils.GenerateSfIDString()
		ext := strings.TrimPrefix(filepath.Ext(param.FileName), ".")
		spaceFile := &driveModel.SpaceFile{
			ID:                  id,
			FileName:            param.FileName,
			FileType:            ext,
			FileCategory:        utils.FileCategoryFromExt(ext),
			FileSize:            param.FileSize,
			PhysicalID:          physicalFile.ID,
			PhysicalStoragePath: physicalFile.StoragePath,
			SpaceID:             space.ID,
			UserID:              userId,
			Path:                path + "/" + id,
			Level:               level + 1,
			ParentID:            parentId,
		}
		// 创建存储关系
		err := driveDao.SpaceFileDao.Create(tx, spaceFile)
		if err != nil {
			return err
		}
		// 物理文件引数加1
		err = driveDao.PhysicalFileDao.FileRefIncById(tx, spaceFile.ID)
		if err != nil {
			return err
		}
		return nil
	})
	return err
}

func (s spaceFileService) VerifySpacePermission(spaceId, userId string) (*driveModel.Space, error) {
	space := driveDao.SpaceDao.GetById(db.RDB, spaceId)
	if space == nil {
		return nil, errors.New("param.error")
	}
	switch space.SpaceType {
	case constant.SpaceType.User:
		if space.OwnerID == userId || space.TargetID == userId {
			return space, nil
		}
	case constant.SpaceType.Group:
		if basicService.GroupService.IsGroupMember(space.TargetID, userId) {
			return space, nil
		}
	}
	return nil, errors.New("param.error")
}

func (s spaceFileService) VerifyDirPermission(parentID, spaceId, userId string) (*driveModel.SpaceFile, error) {
	if parentID == "root" {
		return nil, nil
	}
	file := driveDao.SpaceFileDao.GetById(db.RDB, parentID)
	if file == nil || file.SpaceID != spaceId || !file.IsDir {
		return nil, errors.New("param.error")
	}
	return file, nil
}
