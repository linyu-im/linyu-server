package dao

import (
	"errors"

	basicModel "github.com/linyu-im/linyu-server/linyu-basic-service/pkg/model"
	"gorm.io/gorm"
)

var AppVersionDao = newAppVersionDao()

func newAppVersionDao() *appVersionDao {
	return &appVersionDao{}
}

type appVersionDao struct{}

func (d *appVersionDao) GetByPlatform(db *gorm.DB, platform string) (*basicModel.AppVersion, error) {
	result := &basicModel.AppVersion{}
	err := db.Where("platform = ?", platform).First(result).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return result, nil
}
