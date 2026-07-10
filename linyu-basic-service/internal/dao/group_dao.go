package dao

import (
	basicModel "github.com/linyu-im/linyu-server/linyu-basic-service/pkg/model"
	"gorm.io/gorm"
)

var GroupDao = newGroupDao()

func newGroupDao() *groupDao {
	return &groupDao{}
}

type groupDao struct{}

func (d *groupDao) GetMaxGroupNumber(db *gorm.DB) (string, error) {
	var maxNumber string
	err := db.Model(&basicModel.Group{}).
		Select("COALESCE(MAX(CAST(group_number AS UNSIGNED)), 0)").
		Scan(&maxNumber).Error
	if err != nil {
		return "", err
	}
	return maxNumber, nil
}

func (d *groupDao) GetGroupByGroupNumber(db *gorm.DB, account string) *basicModel.Group {
	result := &basicModel.Group{}
	if err := db.First(result, "group_number = ?", account).Error; err != nil {
		return nil
	}
	return result
}

func (d *groupDao) Create(db *gorm.DB, group *basicModel.Group) error {
	if err := db.Create(group).Error; err != nil {
		return err
	}
	return nil
}

func (d *groupDao) DeleteById(db *gorm.DB, id string) error {
	err := db.Where("id = ?", id).Delete(&basicModel.Group{}).Error
	return err
}

func (d *groupDao) GetGroupById(db *gorm.DB, id string) *basicModel.Group {
	result := &basicModel.Group{}
	if err := db.First(result, "id = ?", id).Error; err != nil {
		return nil
	}
	return result
}

func (d *groupDao) UpdateMemberNum(db *gorm.DB, groupId string) error {
	var memberCount int64
	if err := db.Model(&basicModel.GroupMember{}).
		Where("group_id = ?", groupId).
		Count(&memberCount).Error; err != nil {
		return err
	}
	if err := db.Model(&basicModel.Group{}).
		Where("id = ?", groupId).
		Update("member_num", memberCount).Error; err != nil {
		return err
	}
	return nil
}

func (d *groupDao) SearchByKeyword(db *gorm.DB, keyword string) *gorm.DB {
	like := "%" + keyword + "%"
	return db.Model(&basicModel.Group{}).Where("group_number LIKE ? OR name LIKE ?", like, like)
}

func (d *groupDao) UpdateAvatar(db *gorm.DB, groupId string, avatar string) error {
	return db.Model(&basicModel.Group{}).
		Where("id = ?", groupId).
		Update("avatar", avatar).Error
}

func (d *groupDao) UpdateInfo(db *gorm.DB, groupId string, fields map[string]interface{}) error {
	return db.Model(&basicModel.Group{}).
		Where("id = ?", groupId).
		Updates(fields).Error
}

func (d *groupDao) GroupInfoById(db *gorm.DB, userId string, groupId string) *basicModel.Group {
	result := &basicModel.Group{}
	if err := db.Table("t_group AS g").
		Select("g.*, gm.group_nick_name, gm.group_remark, gm.group_user_level").
		Joins(
			"LEFT JOIN t_group_member AS gm ON gm.group_id = g.id AND gm.user_id = ?",
			userId,
		).
		Where("g.id = ?", groupId).
		First(result).Error; err != nil {
		return nil
	}
	return result
}
