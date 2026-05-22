package dao

import (
	basicModel "github.com/linyu-im/linyu-server/linyu-basic-service/pkg/model"
	"gorm.io/gorm"
)

var MomentCommentDao = newMomentCommentDao()

func newMomentCommentDao() *momentCommentDao {
	return &momentCommentDao{}
}

type momentCommentDao struct{}

func (d momentCommentDao) ListByMomentId(db *gorm.DB, momentId string) []*basicModel.MomentComment {
	var list []*basicModel.MomentComment
	db.Table("t_moment_comment c").
		Select(`
			c.*,
			u.username AS username,
			u.avatar AS user_avatar
		`).
		Joins("LEFT JOIN t_user u ON u.id = c.user_id").
		Where("c.moment_id = ?", momentId).
		Order("c.created_at ASC").
		Scan(&list)
	return list
}

func (d momentCommentDao) Create(db *gorm.DB, comment *basicModel.MomentComment) error {
	if err := db.Create(comment).Error; err != nil {
		return err
	}
	return nil
}

func (d momentCommentDao) GetById(db *gorm.DB, id string) *basicModel.MomentComment {
	result := &basicModel.MomentComment{}
	if err := db.Table("t_moment_comment c").
		Select(`
			c.*,
			u.username AS username,
			u.avatar AS user_avatar
		`).
		Joins("LEFT JOIN t_user u ON u.id = c.user_id").
		Where("c.id = ?", id).
		First(result).Error; err != nil {
		return nil
	}
	return result
}

func (d momentCommentDao) DeleteByUserIdAndMomentId(db *gorm.DB, userId string, commentId string) error {
	return db.Delete(&basicModel.MomentComment{}, "user_id = ? AND id = ?", userId, commentId).Error
}
