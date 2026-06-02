package service

import (
	"errors"

	basicDao "github.com/linyu-im/linyu-server/linyu-basic-service/internal/dao"
	basicModel "github.com/linyu-im/linyu-server/linyu-basic-service/pkg/model"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/db"
)

var EnterpriseService = newEnterpriseService()

func newEnterpriseService() *enterpriseService {
	return &enterpriseService{}
}

type enterpriseService struct{}

func (s *enterpriseService) EnterpriseInfo(userId string, enterpriseId string) (*basicModel.Enterprise, error) {
	if basicDao.EnterpriseMemberDao.GetByEnterpriseIdAndUserId(db.RDB, enterpriseId, userId) == nil {
		return nil, errors.New("param.error")
	}
	enterprise := basicDao.EnterpriseDao.GetById(db.RDB, enterpriseId)
	if enterprise == nil {
		return nil, errors.New("param.error")
	}

	userMembers, err := basicDao.EnterpriseMemberDao.ListByEnterpriseIdAndUserId(db.RDB, enterpriseId, userId)
	if err != nil {
		return nil, err
	}
	enterprise.UserEnterpriseMembers = userMembers
	return enterprise, nil
}

func (s *enterpriseService) EnterpriseDepartment(enterpriseId string, userId string) ([]*basicModel.EnterpriseDepartment, error) {
	if basicDao.EnterpriseMemberDao.GetByEnterpriseIdAndUserId(db.RDB, enterpriseId, userId) == nil {
		return nil, errors.New("param.error")
	}
	departments, err := basicDao.EnterpriseDepartmentDao.ListByEnterpriseId(db.RDB, enterpriseId)
	if err != nil {
		return nil, err
	}
	members, err := basicDao.EnterpriseMemberDao.ListByEnterpriseId(db.RDB, enterpriseId)
	if err != nil {
		return nil, err
	}
	return buildEnterpriseDepartmentTree(departments, members), nil
}

func buildEnterpriseDepartmentTree(
	departments []*basicModel.EnterpriseDepartment,
	members []*basicModel.EnterpriseMember,
) []*basicModel.EnterpriseDepartment {
	deptMap := make(map[string]*basicModel.EnterpriseDepartment, len(departments))
	for _, dept := range departments {
		dept.Children = []*basicModel.EnterpriseDepartment{}
		dept.Members = []*basicModel.EnterpriseMember{}
		deptMap[dept.ID] = dept
	}
	for _, member := range members {
		if member.DepartmentID == "" {
			continue
		}
		if dept, ok := deptMap[member.DepartmentID]; ok {
			dept.Members = append(dept.Members, member)
		}
	}
	roots := make([]*basicModel.EnterpriseDepartment, 0)
	for _, dept := range departments {
		if dept.ParentID == "" {
			roots = append(roots, dept)
			continue
		}
		if parent, ok := deptMap[dept.ParentID]; ok {
			parent.Children = append(parent.Children, dept)
		}
	}
	return roots
}

func (s *enterpriseService) IsEnterpriseMember(enterpriseId string, userId string) bool {
	return basicDao.EnterpriseMemberDao.GetByEnterpriseIdAndUserId(db.RDB, enterpriseId, userId) != nil
}

func (s *enterpriseService) GetEnterpriseAvatar(enterpriseId string) interface{} {
	enterprise := basicDao.EnterpriseDao.GetById(db.RDB, enterpriseId)
	if enterprise == nil {
		return ""
	}
	return enterprise.Avatar
}
