package dao

import (
	basicModel "github.com/linyu-im/linyu-server/linyu-basic-service/pkg/model"
	"gorm.io/gorm"
)

var GroupDao = newGroupDao()

func newGroupDao() *groupDao {
	return &groupDao{}
}

type groupDao struct{}

func (d *groupDao) GetGroupByGroupNumber(db *gorm.DB, account string) *basicModel.Group {
	result := &basicModel.Group{}
	if err := db.First(result, "group_number = ?", account).Error; err != nil {
		return nil
	}
	return result
}

func (d *groupDao) Create(db *gorm.DB, group *basicModel.Group) error {
	if err := db.Create(group).Error; err != nil {
		return err
	}
	return nil
}

func (d *groupDao) DeleteById(db *gorm.DB, id string) error {
	err := db.Where("id = ?", id).Delete(&basicModel.Group{}).Error
	return err
}

func (d *groupDao) GetGroupById(db *gorm.DB, id string) *basicModel.Group {
	result := &basicModel.Group{}
	if err := db.First(result, "id = ?", id).Error; err != nil {
		return nil
	}
	return result
}
