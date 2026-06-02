package dao

import (
	basicModel "github.com/linyu-im/linyu-server/linyu-basic-service/pkg/model"
	"gorm.io/gorm"
)

var EnterpriseDao = newEnterpriseDao()

func newEnterpriseDao() *enterpriseDao {
	return &enterpriseDao{}
}

type enterpriseDao struct{}

func (d *enterpriseDao) GetById(db *gorm.DB, id string) *basicModel.Enterprise {
	result := &basicModel.Enterprise{}
	if err := db.First(result, "id = ?", id).Error; err != nil {
		return nil
	}
	return result
}
