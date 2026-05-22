package dao

import (
	basicModel "github.com/linyu-im/linyu-server/linyu-basic-service/pkg/model"
	"gorm.io/gorm"
)

var MomentLikeDao = newMomentLikeDao()

func newMomentLikeDao() *momentLikeDao {
	return &momentLikeDao{}
}

type momentLikeDao struct{}

func (d momentLikeDao) Create(db *gorm.DB, like *basicModel.MomentLike) error {
	if err := db.Create(like).Error; err != nil {
		return err
	}
	return nil
}

func (d momentLikeDao) ListByMomentId(db *gorm.DB, momentId string) []*basicModel.MomentLike {
	var list []*basicModel.MomentLike
	db.Table("t_moment_like l").
		Select(`
			l.*,
			u.username AS username,
			u.avatar AS user_avatar
		`).
		Joins("LEFT JOIN t_user u ON u.id = l.user_id").
		Where("l.moment_id = ?", momentId).
		Order("l.created_at ASC").
		Scan(&list)
	return list
}

func (d momentLikeDao) DeleteByUserIdAndMomentId(db *gorm.DB, userId string, momentId string) error {
	return db.Delete(&basicModel.MomentLike{}, "user_id = ? AND moment_id = ?", userId, momentId).Error
}

func (d momentLikeDao) GetByUserIdAndMomentId(db *gorm.DB, userId string, momentId string) *basicModel.MomentLike {
	result := &basicModel.MomentLike{}
	if err := db.First(result, "user_id = ? AND moment_id = ?", userId, momentId).Error; err != nil {
		return nil
	}
	return result
}
