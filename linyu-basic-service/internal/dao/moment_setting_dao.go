package dao

import (
	basicModel "github.com/linyu-im/linyu-server/linyu-basic-service/pkg/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var MomentSetDao = newMomentSetDao()

func newMomentSetDao() *momentSetDao {
	return &momentSetDao{}
}

type momentSetDao struct{}

func (d *momentSetDao) UpdateBackground(db *gorm.DB, userId string, bgUrl string) error {
	setting := &basicModel.MomentSetting{
		UserID: userId,
		BgUrl:  bgUrl,
	}
	return db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"bg_url", "updated_at"}),
	}).Create(setting).Error
}

func (d *momentSetDao) UpdateExpireDays(db *gorm.DB, userId string, expireDays int64) error {
	setting := &basicModel.MomentSetting{
		UserID:     userId,
		ExpireDays: expireDays,
	}
	return db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"expire_days", "updated_at"}),
	}).Create(setting).Error
}

func (d *momentSetDao) GetByUserId(db *gorm.DB, userId string) *basicModel.MomentSetting {
	result := &basicModel.MomentSetting{}
	if err := db.Where("user_id = ?", userId).First(result).Error; err != nil {
		return nil
	}
	return result
}

func (d *momentSetDao) Create(db *gorm.DB, setting *basicModel.MomentSetting) error {
	return db.Create(setting).Error
}
