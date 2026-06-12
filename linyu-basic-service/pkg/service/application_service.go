package service

import (
	basicDao "github.com/linyu-im/linyu-server/linyu-basic-service/internal/dao"
	basicModel "github.com/linyu-im/linyu-server/linyu-basic-service/pkg/model"
	basicParam "github.com/linyu-im/linyu-server/linyu-basic-service/pkg/param"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/db"
)

var ApplicationService = newApplicationService()

func newApplicationService() *applicationService {
	return &applicationService{}
}

type applicationService struct{}

func (s applicationService) ApplicationList(param *basicParam.ApplicationListParam) ([]*basicModel.Application, error) {
	return basicDao.ApplicationDao.List(db.RDB, param.Keyword)
}
