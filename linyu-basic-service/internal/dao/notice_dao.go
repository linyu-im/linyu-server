package dao

import (
	basicModel "github.com/linyu-im/linyu-server/linyu-basic-service/pkg/model"
	"gorm.io/gorm"
)

var NoticeDao = newNoticeDao()

func newNoticeDao() *noticeDao {
	return &noticeDao{}
}

type noticeDao struct{}

func (d *noticeDao) CreateBatch(db *gorm.DB, notices []*basicModel.Notice) error {
	if len(notices) == 0 {
		return nil
	}
	return db.Create(&notices).Error
}

func (d *noticeDao) ListByUserIdAndType(db *gorm.DB, userId string, noticeType string) ([]*basicModel.Notice, error) {
	var list []*basicModel.Notice
	if err := db.Where("user_id = ? AND type = ?", userId, noticeType).
		Order("created_at DESC").
		Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (d *noticeDao) CountAfterLastReadId(db *gorm.DB, userId string, noticeType string, lastReadId string) (int64, error) {
	tx := db.Model(&basicModel.Notice{}).
		Where("user_id = ? AND type = ?", userId, noticeType)
	if lastReadId != "" && lastReadId != "0" {
		tx = tx.Where("id > ?", lastReadId)
	}
	var count int64
	if err := tx.Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
