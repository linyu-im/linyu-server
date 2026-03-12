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
