package dao

import (
	basicModel "github.com/linyu-im/linyu-server/linyu-basic-service/pkg/model"
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

func (d *applyDao) ApplyListAndPeer(db *gorm.DB, userId string, applyType string) ([]*basicModel.Apply, error) {
	var applyList []*basicModel.Apply
	if err := db.Table("t_apply").
		Select("t_apply.*, t_user.username AS peer_name").
		Joins("LEFT JOIN t_user ON t_apply.peer_id = t_user.id").
		Where("(user_id = ? OR peer_id = ?) AND type = ?", userId, userId, applyType).
		Find(&applyList).Error; err != nil {
		return nil, err
	}
	return applyList, nil
}
