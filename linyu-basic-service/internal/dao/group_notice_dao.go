package dao

import (
	basicModel "github.com/linyu-im/linyu-server/linyu-basic-service/pkg/model"
	"gorm.io/gorm"
)

var GroupNoticeDao = newGroupNoticeDao()

func newGroupNoticeDao() *groupNoticeDao {
	return &groupNoticeDao{}
}

type groupNoticeDao struct{}

func (d *groupNoticeDao) Create(db *gorm.DB, notice *basicModel.GroupNotice) error {
	return db.Create(notice).Error
}

func (d *groupNoticeDao) GetById(db *gorm.DB, noticeId string) *basicModel.GroupNotice {
	result := &basicModel.GroupNotice{}
	if err := db.First(result, "id = ?", noticeId).Error; err != nil {
		return nil
	}
	return result
}

func (d *groupNoticeDao) Update(db *gorm.DB, noticeId string, fields map[string]interface{}) error {
	return db.Model(&basicModel.GroupNotice{}).
		Where("id = ?", noticeId).
		Updates(fields).Error
}

func (d *groupNoticeDao) DeleteById(db *gorm.DB, noticeId string) error {
	return db.Where("id = ?", noticeId).Delete(&basicModel.GroupNotice{}).Error
}

func (d *groupNoticeDao) ListByGroupId(db *gorm.DB, groupId string) ([]*basicModel.GroupNotice, error) {
	var notices []*basicModel.GroupNotice
	if err := db.Where("group_id = ?", groupId).
		Order("is_top DESC, created_at DESC").
		Find(&notices).Error; err != nil {
		return nil, err
	}
	return notices, nil
}

func (d *groupNoticeDao) ResetIsTopByGroupId(db *gorm.DB, groupId string) error {
	return db.Model(&basicModel.GroupNotice{}).
		Where("group_id = ? AND is_top = ?", groupId, true).
		Update("is_top", false).Error
}

func (d *groupNoticeDao) GetLatestContentByGroupId(db *gorm.DB, groupId string) string {
	notice := &basicModel.GroupNotice{}
	if err := db.Where("group_id = ? AND is_top = ?", groupId, true).First(notice).Error; err == nil {
		return notice.Content
	}
	if err := db.Where("group_id = ?", groupId).Order("created_at DESC").First(notice).Error; err == nil {
		return notice.Content
	}
	return ""
}
