package dao

import (
	basicModel "github.com/linyu-im/linyu-server/linyu-basic-service/pkg/model"
	"gorm.io/gorm"
)

var MomentVisibleDao = newMomentVisibleDao()

func newMomentVisibleDao() *momentVisibleDao {
	return &momentVisibleDao{}
}

type momentVisibleDao struct{}

func (d momentVisibleDao) CreateBatch(db *gorm.DB, visibles []*basicModel.MomentVisible) error {
	if len(visibles) == 0 {
		return nil
	}
	err := db.CreateInBatches(visibles, 500).Error
	if err != nil {
		return err
	}
	return nil
}
