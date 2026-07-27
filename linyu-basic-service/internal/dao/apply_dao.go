package dao

import (
	basicModel "github.com/linyu-im/linyu-server/linyu-basic-service/pkg/model"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/constant"
	"gorm.io/gorm"
)

var ApplyDao = newApplyDao()

func newApplyDao() *applyDao {
	return &applyDao{}
}

type applyDao struct{}

func (d *applyDao) Create(db *gorm.DB, apply *basicModel.Apply) error {
	if err := db.Create(apply).Error; err != nil {
		return err
	}
	return nil
}

func (d *applyDao) GetById(db *gorm.DB, applyId string) *basicModel.Apply {
	result := &basicModel.Apply{}
	if err := db.First(result, "id = ?", applyId).Error; err != nil {
		return nil
	}
	return result
}

func (d *applyDao) Update(db *gorm.DB, apply *basicModel.Apply) error {
	if err := db.Updates(apply).Error; err != nil {
		return err
	}
	return nil
}

func (d *applyDao) ApplyFriendList(db *gorm.DB, userId string) ([]*basicModel.Apply, error) {
	var applyList []*basicModel.Apply
	if err := db.Table("t_apply").
		Select("t_apply.*, t_user.username AS peer_name").
		Joins("LEFT JOIN t_user ON t_apply.peer_id = t_user.id").
		Where("(user_id = ? OR peer_id = ?) AND type = ?", userId, userId, constant.ApplyType.Friend).
		Order("t_apply.created_at DESC").
		Find(&applyList).Error; err != nil {
		return nil, err
	}
	return applyList, nil
}

func (d *applyDao) ApplyGroupList(db *gorm.DB, userId string) ([]*basicModel.Apply, error) {
	var applyList []*basicModel.Apply
	err := db.Table("t_apply AS a").
		Select("a.*, g.name AS peer_name, u.username AS user_name").
		Joins("LEFT JOIN t_group g ON a.peer_id = g.id AND g.deleted_at IS NULL").
		Joins("LEFT JOIN t_user u ON a.user_id = u.id AND u.deleted_at IS NULL").
		Where(
			`a.type = ? AND a.deleted_at IS NULL AND (
				a.user_id = ?
				OR a.peer_id IN (
					SELECT gm.group_id
					FROM t_group_member gm
					WHERE gm.user_id = ?
					AND gm.member_role = ?
					AND gm.deleted_at IS NULL
				)
			)`,
			constant.ApplyType.Group,
			userId,
			userId,
			constant.MemberRole.Admin,
		).
		Order("a.created_at DESC").
		Find(&applyList).Error
	if err != nil {
		return nil, err
	}
	return applyList, nil
}

func (d *applyDao) CountFriendApplyAfterLastReadId(db *gorm.DB, userId string, lastReadId string) (int64, error) {
	tx := db.Table("t_apply").
		Where("(user_id = ? OR peer_id = ?) AND type = ?", userId, userId, constant.ApplyType.Friend).
		Where("NOT (user_id = ? AND status = ?)", userId, constant.ApplyStatus.Wait)
	if lastReadId != "" && lastReadId != "0" {
		tx = tx.Where("id > ?", lastReadId)
	}
	var count int64
	if err := tx.Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
