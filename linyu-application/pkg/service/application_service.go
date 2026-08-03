package service

import (
	appDao "github.com/linyu-im/linyu-server/linyu-application/internal/dao"
	appModel "github.com/linyu-im/linyu-server/linyu-application/pkg/model"
	appParam "github.com/linyu-im/linyu-server/linyu-application/pkg/param"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/db"
)

var ApplicationService = newApplicationService()

func newApplicationService() *applicationService {
	return &applicationService{}
}

type applicationService struct{}

func (s applicationService) ApplicationList(param *appParam.ApplicationListParam) ([]*appModel.Application, error) {
	return appDao.ApplicationDao.List(db.RDB, param.Keyword)
}
