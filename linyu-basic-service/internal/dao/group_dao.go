package dao

import (
	"sync"

	basicModel "github.com/linyu-im/linyu-server/linyu-basic-service/pkg/model"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/constant"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/db"
	"gorm.io/gorm"
)

var GroupDao = newGroupDao()

func newGroupDao() *groupDao {
	return &groupDao{}
}

type groupDao struct{}

var maxGroupNumberMu sync.Mutex

func (d *groupDao) GetMaxGroupNumber(gdb *gorm.DB) (string, error) {
	key := constant.RedisKey.GroupMaxNumber
	if val, err := db.CacheDB.Get(key); err == nil && val != "" {
		return val, nil
	}

	// 缓存未命中时加锁，避免并发打穿数据库
	maxGroupNumberMu.Lock()
	defer maxGroupNumberMu.Unlock()

	if val, err := db.CacheDB.Get(key); err == nil && val != "" {
		return val, nil
	}

	var maxNumber string
	err := gdb.Model(&basicModel.Group{}).
		Select("COALESCE(MAX(CAST(group_number AS UNSIGNED)), 0)").
		Scan(&maxNumber).Error
	if err != nil {
		return "", err
	}
	_ = db.CacheDB.Set(key, maxNumber, 0)
	return maxNumber, nil
}

// UpdateMaxGroupNumberCache 创建群成功后回写最大群号缓存（原子取较大值）
func (d *groupDao) UpdateMaxGroupNumberCache(number string) {
	if number == "" {
		return
	}
	_ = db.CacheDB.SetIfGreater(constant.RedisKey.GroupMaxNumber, number, 0)
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

func (d *groupDao) UnscopedDeleteById(db *gorm.DB, id string) error {
	err := db.Unscoped().Where("id = ?", id).Delete(&basicModel.Group{}).Error
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

func (d *groupDao) UpdateOwnerUserId(db *gorm.DB, groupId string, newOwnerId string) error {
	return db.Model(&basicModel.Group{}).
		Where("id = ?", groupId).
		Update("owner_user_id", newOwnerId).Error
}

func (d *groupDao) GroupInfoById(db *gorm.DB, userId string, groupId string) *basicModel.Group {
	result := &basicModel.Group{}
	if err := db.Unscoped().Table("t_group AS g").
		Select("g.*, gm.group_nick_name, gm.group_user_level").
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
