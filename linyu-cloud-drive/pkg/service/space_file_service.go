package service

import (
	"errors"
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

func (s spaceFileService) CreateFileFromPhysicalFile(userId string,
	param *driveParam.UploadFileInfoParam,
	physicalFile *driveModel.PhysicalFile) error {

	if !constant.SpaceType.Validate(param.SpaceType) {
		return errors.New("param.error")
	}
	// 验证目录权限
	parentFile, err := s.VerifyDirPermission(param.ParentID, userId)
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
		spaceFile := &driveModel.SpaceFile{
			ID:                  utils.GenerateSfIDString(),
			FileName:            param.FileName,
			FileSize:            param.FileSize,
			PhysicalID:          physicalFile.ID,
			PhysicalStoragePath: physicalFile.StoragePath,
			SpaceID:             param.SpaceID,
			SpaceType:           param.SpaceType,
			UserID:              userId,
			Path:                path + "/" + param.FileName,
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

func (s spaceFileService) VerifyDirPermission(parentID, userId string) (*driveModel.SpaceFile, error) {
	// 验证父目录id是否是文件夹或者根目录
	if parentID != "root" {
		file := driveDao.SpaceFileDao.GetById(db.RDB, parentID)
		if file == nil || !file.IsDir {
			return nil, errors.New("param.error")
		}
		// 验证不同类型的操作权限
		switch file.SpaceType {
		case constant.SpaceType.User:
			if file.SpaceID == userId {
				return file, nil
			}
		case constant.SpaceType.Group:
			is := basicService.GroupService.IsGroupMember(file.SpaceID, userId)
			if is {
				return file, nil
			}
		}
	}
	return nil, nil
}
