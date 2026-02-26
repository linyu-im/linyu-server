package service

import (
	driveDao "github.com/linyu-im/linyu-server/linyu-cloud-drive/internal/dao"
	driveModel "github.com/linyu-im/linyu-server/linyu-cloud-drive/pkg/model"
	driveParam "github.com/linyu-im/linyu-server/linyu-cloud-drive/pkg/param"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/db"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/utils"
)

var PhysicalFileService = newPhysicalFileService()

func newPhysicalFileService() *physicalFileService {
	return &physicalFileService{}
}

type physicalFileService struct{}

func (s physicalFileService) GetFileByHash(hash string) *driveModel.PhysicalFile {
	return driveDao.PhysicalFileDao.GetByHash(db.RDB, hash)
}

func (s physicalFileService) CreateFile(param *driveParam.UploadFileInfoParam, storagePath string) (error, *driveModel.PhysicalFile) {
	file := &driveModel.PhysicalFile{
		ID:          utils.GenerateSfIDString(),
		FileHash:    param.FileHash,
		FileSize:    param.FileSize,
		StoragePath: storagePath,
	}
	err := driveDao.PhysicalFileDao.Create(db.RDB, file)
	if err != nil {
		return err, nil
	}
	return nil, file
}
