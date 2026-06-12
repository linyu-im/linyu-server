package dao

import (
	basicModel "github.com/linyu-im/linyu-server/linyu-basic-service/pkg/model"
	"gorm.io/gorm"
)

var ApplicationDao = newApplicationDao()

func newApplicationDao() *applicationDao {
	return &applicationDao{}
}

type applicationDao struct{}

func (d *applicationDao) List(db *gorm.DB, keyword string) ([]*basicModel.Application, error) {
	var list []*basicModel.Application
	tx := db.Model(&basicModel.Application{})
	if keyword != "" {
		like := "%" + keyword + "%"
		tx = tx.Where("app_name LIKE ? OR author LIKE ? OR description LIKE ?", like, like, like)
	}
	if err := tx.Order("created_at DESC").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}
