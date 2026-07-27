package dao

import (
	basicModel "github.com/linyu-im/linyu-server/linyu-basic-service/pkg/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var UserBadgeDao = newUserBadgeDao()

func newUserBadgeDao() *userBadgeDao {
	return &userBadgeDao{}
}

type userBadgeDao struct{}

func (d *userBadgeDao) UpsertLastReadID(db *gorm.DB, userId string, badgeCode string, lastReadId string) error {
	badge := &basicModel.UserBadge{
		UserID:     userId,
		BadgeCode:  badgeCode,
		LastReadID: lastReadId,
	}
	return db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}, {Name: "badge_code"}},
		DoUpdates: clause.AssignmentColumns([]string{"last_read_id", "updated_at"}),
	}).Create(badge).Error
}

func (d *userBadgeDao) ListByUserId(db *gorm.DB, userId string) ([]*basicModel.UserBadge, error) {
	var list []*basicModel.UserBadge
	if err := db.Where("user_id = ?", userId).Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (d *userBadgeDao) GetByUserIdAndCode(db *gorm.DB, userId string, badgeCode string) *basicModel.UserBadge {
	result := &basicModel.UserBadge{}
	if err := db.Where("user_id = ? AND badge_code = ?", userId, badgeCode).First(result).Error; err != nil {
		return nil
	}
	return result
}

func (d *userBadgeDao) Create(db *gorm.DB, badge *basicModel.UserBadge) error {
	return db.Create(badge).Error
}
