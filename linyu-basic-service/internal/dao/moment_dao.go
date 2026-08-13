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

	// 查看自己全部（不过滤可见性和过期）
	if viewUserId == userId {
		return tx.Where("t_moment.user_id = ?", userId).
			Order("t_moment.created_at DESC")
	}

	// 查看指定好友
	if viewUserId != "" {
		// 过期规则：NULL/0=永久；>0=最近 N 天；-1=对他人不可见
		now := time.Now()
		tx = tx.Where(`(
				ms.expire_days IS NULL
				OR ms.expire_days = 0
				OR (
					ms.expire_days > 0
					AND t_moment.created_at >= DATE_SUB(?, INTERVAL ms.expire_days DAY)
				)
			)`, now)
		tx = tx.Where("t_moment.user_id = ?", viewUserId)
		tx = d.buildVisibilityWhere(tx, userId)
		return tx.Order("t_moment.created_at DESC")
	}

	// 查看全部：自己的 + 好友的（自己的不受可见性和过期限制）
	now := time.Now()
	tx = tx.Where(`(
		t_moment.user_id = ?
		OR (
			t_moment.user_id IN (
				SELECT peer_id FROM t_contacts WHERE user_id = ?
			)
			AND (
				ms.expire_days IS NULL
				OR ms.expire_days = 0
				OR (
					ms.expire_days > 0
					AND t_moment.created_at >= DATE_SUB(?, INTERVAL ms.expire_days DAY)
				)
			)
		)
	)`, userId, userId, now)
	tx = d.buildVisibilityWhere(tx, userId)

	return tx.Order("t_moment.created_at DESC")
}

func (d momentDao) DeleteByUserIdAndMomentId(db *gorm.DB, userId string, momentId string) error {
	return db.Delete(&basicModel.Moment{}, "user_id = ? AND id = ?", userId, momentId).Error
}

func (d momentDao) buildVisibilityWhere(tx *gorm.DB, userId string) *gorm.DB {
	return tx.Where(`
		(
			t_moment.user_id = ?
			OR
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
	`, userId, userId, userId)
}
