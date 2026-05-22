package dao

import (
	basicModel "github.com/linyu-im/linyu-server/linyu-basic-service/pkg/model"
	"gorm.io/gorm"
	"time"
)

var MomentDao = newMomentDao()

func newMomentDao() *momentDao {
	return &momentDao{}
}

type momentDao struct{}

func (d momentDao) Create(db *gorm.DB, moment *basicModel.Moment) error {
	if err := db.Create(moment).Error; err != nil {
		return err
	}
	return nil
}

func (d momentDao) GetMomentById(db *gorm.DB, momentId string) *basicModel.Moment {
	result := &basicModel.Moment{}
	if err := db.First(result, "id = ?", momentId).Error; err != nil {
		return nil
	}
	return result
}

func (d momentDao) BuildMomentQuery(db *gorm.DB, userId string, viewUserId string) *gorm.DB {

	tx := db.Table("t_moment").
		Select("t_moment.*, u.username AS username, u.user_level AS userLevel").
		Joins("LEFT JOIN t_moment_setting ms ON ms.user_id = t_moment.user_id").
		Joins("LEFT JOIN t_user u ON u.id = t_moment.user_id AND u.deleted_at IS NULL").
		Where("t_moment.deleted_at IS NULL")

	// 查看自己
	if viewUserId == userId {
		tx = tx.Where("t_moment.user_id = ?", userId).
			Order("t_moment.created_at DESC")
		return tx
	}

	// 过期规则
	now := time.Now()
	tx = tx.Where(`
			ms.expire_days IS NULL
			OR ms.expire_days = 0
			OR (
				ms.expire_days > 0
				AND t_moment.created_at >= DATE_SUB(?, INTERVAL ms.expire_days DAY)
			)
		`, now)

	// 查看指定好友
	if viewUserId != "" {

		tx = tx.Where("t_moment.user_id = ?", viewUserId)

	} else {

		// 查看所有好友朋友圈
		tx = tx.Where(`
			t_moment.user_id IN (
				SELECT peer_id
				FROM t_contacts
				WHERE user_id = ?
			)
		`, userId)

	}

	// 可见性规则
	tx = tx.Where(`
		(
			t_moment.visible_type = 'all'
			OR
			(
				t_moment.visible_type = 'include'
				AND EXISTS (
					SELECT 1 FROM t_moment_visible mv
					WHERE mv.moment_id = t_moment.id
					AND mv.user_id = ?
					AND mv.visible_type = 'include'
				)
			)
			OR
			(
				t_moment.visible_type = 'exclude'
				AND NOT EXISTS (
					SELECT 1 FROM t_moment_visible mv
					WHERE mv.moment_id = t_moment.id
					AND mv.user_id = ?
					AND mv.visible_type = 'exclude'
				)
			)
		)
	`, userId, userId)

	return tx.Order("t_moment.created_at DESC")
}
