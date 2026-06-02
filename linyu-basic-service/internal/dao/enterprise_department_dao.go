package dao

import (
	basicModel "github.com/linyu-im/linyu-server/linyu-basic-service/pkg/model"
	"gorm.io/gorm"
)

var EnterpriseDepartmentDao = newEnterpriseDepartmentDao()

func newEnterpriseDepartmentDao() *enterpriseDepartmentDao {
	return &enterpriseDepartmentDao{}
}

type enterpriseDepartmentDao struct{}

func (d *enterpriseDepartmentDao) ListByEnterpriseId(db *gorm.DB, enterpriseId string) ([]*basicModel.EnterpriseDepartment, error) {
	var list []*basicModel.EnterpriseDepartment
	err := db.Table("t_enterprise_department AS d").
		Select("d.*, u.username AS leader_username").
		Joins("LEFT JOIN t_user u ON u.id = d.leader_user_id AND u.deleted_at IS NULL").
		Where("d.enterprise_id = ? AND d.deleted_at IS NULL", enterpriseId).
		Order("d.sort ASC, d.level ASC, d.created_at ASC").
		Find(&list).Error
	if err != nil {
		return nil, err
	}
	return list, nil
}
