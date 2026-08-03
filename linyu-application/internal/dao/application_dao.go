package dao

import (
	appModel "github.com/linyu-im/linyu-server/linyu-application/pkg/model"
	"gorm.io/gorm"
)

var ApplicationDao = newApplicationDao()

func newApplicationDao() *applicationDao {
	return &applicationDao{}
}

type applicationDao struct{}

func (d *applicationDao) List(db *gorm.DB, keyword string) ([]*appModel.Application, error) {
	var list []*appModel.Application
	tx := db.Model(&appModel.Application{})
	if keyword != "" {
		like := "%" + keyword + "%"
		tx = tx.Where("app_name LIKE ? OR author LIKE ? OR description LIKE ?", like, like, like)
	}
	if err := tx.Order("created_at DESC").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}
