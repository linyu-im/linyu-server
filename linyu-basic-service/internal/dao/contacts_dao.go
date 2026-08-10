package dao

import (
	basicModel "github.com/linyu-im/linyu-server/linyu-basic-service/pkg/model"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/constant"
	"gorm.io/gorm"
)

var ContactsDao = newContactsDao()

func newContactsDao() *contactsDao {
	return &contactsDao{}
}

type contactsDao struct{}

func (d *contactsDao) IsContactByUserAndPeer(db *gorm.DB, userId string, peerId string) bool {
	var count int64
	err := db.Model(&basicModel.Contacts{}).
		Where(
			"(user_id = ? AND peer_id = ?) OR (user_id = ? AND peer_id = ?)",
			userId, peerId,
			peerId, userId,
		).
		Count(&count).
		Error
	if err != nil {
		return false
	}
	return count > 0
}

func (d *contactsDao) IsFriend(db *gorm.DB, userId string, peerId string) bool {
	var count int64
	err := db.Model(&basicModel.Contacts{}).
		Where(
			"(user_id = ? AND peer_id = ?)",
			userId, peerId,
		).
		Count(&count).
		Error
	if err != nil {
		return false
	}
	return count > 0
}

func (d *contactsDao) IsFriendBothOr(db *gorm.DB, userId string, peerId string) bool {
	var count int64
	err := db.Model(&basicModel.Contacts{}).
		Where(
			"(user_id = ? AND peer_id = ?) OR (user_id = ? AND peer_id = ?)",
			userId, peerId,
			peerId, userId,
		).
		Count(&count).
		Error
	if err != nil {
		return false
	}
	return count > 0
}

func (d *contactsDao) IsFriendBothAnd(db *gorm.DB, userId string, peerId string) bool {
	return d.IsFriend(db, userId, peerId) && d.IsFriend(db, peerId, userId)
}

func (d *contactsDao) IsGroupContact(db *gorm.DB, userId string, groupId string) bool {
	var count int64
	err := db.Model(&basicModel.Contacts{}).
		Where("user_id = ? AND peer_id = ? AND peer_type = ?", userId, groupId, constant.ContactsPeerType.Group).
		Count(&count).Error
	if err != nil {
		return false
	}
	return count > 0
}

func (d *contactsDao) Create(db *gorm.DB, contacts *basicModel.Contacts) error {
	if err := db.Create(contacts).Error; err != nil {
		return err
	}
	return nil
}

func (d *contactsDao) ContactsFriendList(db *gorm.DB, userId string) ([]*basicModel.Contacts, error) {
	var list []*basicModel.Contacts
	if err := db.Table("t_contacts AS c").
		Select("c.*, u.username AS username, u.user_level AS user_level, e.emotion_name AS emotion_name, e.url AS emotion_url ").
		Joins("LEFT JOIN t_user u ON u.id = c.peer_id").
		Joins("LEFT JOIN t_emotion e ON u.emotion_id = e.id").
		Where("c.user_id = ? AND c.peer_type = ?", userId, constant.ContactsPeerType.Friend).
		Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (d *contactsDao) ContactsFriendSearch(db *gorm.DB, userId string, keyword string) ([]*basicModel.Contacts, error) {
	var list []*basicModel.Contacts
	like := "%" + keyword + "%"
	if err := db.Table("t_contacts AS c").
		Select("c.*, u.username AS username, u.user_level AS user_level, e.emotion_name AS emotion_name, e.url AS emotion_url ").
		Joins("LEFT JOIN t_user u ON u.id = c.peer_id").
		Joins("LEFT JOIN t_emotion e ON u.emotion_id = e.id").
		Where("c.user_id = ? AND c.peer_type = ?", userId, constant.ContactsPeerType.Friend).
		Where("u.username LIKE ? OR u.account LIKE ? OR c.remark LIKE ?", like, like, like).
		Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (d *contactsDao) ContactsGroupList(db *gorm.DB, userId string) ([]*basicModel.Contacts, error) {
	var list []*basicModel.Contacts
	if err := db.Table("t_contacts AS c").
		Select("c.*, g.name AS group_name, g.member_num AS group_member_num").
		Joins("LEFT JOIN t_group g ON g.id = c.peer_id").
		Where("c.user_id = ? AND c.peer_type = ?", userId, constant.ContactsPeerType.Group).
		Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (d *contactsDao) ContactsGroupSearch(db *gorm.DB, userId string, keyword string) ([]*basicModel.Contacts, error) {
	var list []*basicModel.Contacts
	like := "%" + keyword + "%"
	if err := db.Table("t_contacts AS c").
		Select("c.*, g.name AS group_name, g.member_num AS group_member_num").
		Joins("LEFT JOIN t_group g ON g.id = c.peer_id").
		Where("c.user_id = ? AND c.peer_type = ?", userId, constant.ContactsPeerType.Group).
		Where("g.name LIKE ? OR g.group_number LIKE ? OR c.remark LIKE ?", like, like, like).
		Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (d *contactsDao) ContactsEnterpriseList(db *gorm.DB, userId string) ([]*basicModel.Contacts, error) {
	var list []*basicModel.Contacts
	if err := db.Table("t_contacts AS c").
		Select("c.*, e.name AS enterprise_name, e.member_num AS enterprise_member_num").
		Joins("LEFT JOIN t_enterprise e ON e.id = c.peer_id").
		Where("c.user_id = ? AND c.peer_type = ?", userId, constant.ContactsPeerType.Enterprise).
		Order("c.created_at DESC").
		Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (d *contactsDao) GetById(db *gorm.DB, contactsId string) (*basicModel.Contacts, error) {
	result := &basicModel.Contacts{}
	if err := db.First(result, "id = ?", contactsId).Error; err != nil {
		return nil, err
	}
	return result, nil
}

func (d *contactsDao) SetIsTopByUserAndPeerId(db *gorm.DB, isTop bool, userId string, peerId string) error {
	return db.Model(&basicModel.Contacts{}).
		Where("user_id = ? AND peer_id = ?", userId, peerId).
		Update("is_top", isTop).Error
}

func (d *contactsDao) SetIsMuteByUserAndPeerId(db *gorm.DB, isMute bool, userId string, peerId string) error {
	return db.Model(&basicModel.Contacts{}).
		Where("user_id = ? AND peer_id = ?", userId, peerId).
		Update("is_mute", isMute).Error
}

func (d *contactsDao) UpdateRemarkByUserAndPeerId(db *gorm.DB, userId string, peerId string, remark string) error {
	return db.Model(&basicModel.Contacts{}).
		Where("user_id = ? AND peer_id = ?", userId, peerId).
		Update("remark", remark).Error
}

func (d *contactsDao) UpdateTagByUserAndPeerId(db *gorm.DB, userId string, peerId string, tag string) error {
	return db.Model(&basicModel.Contacts{}).
		Where("user_id = ? AND peer_id = ?", userId, peerId).
		Update("tag", tag).Error
}

func (d *contactsDao) UnscopedDeleteByUserAndPeerId(db *gorm.DB, userId string, peerId string) error {
	err := db.Unscoped().Where("user_id = ? AND peer_id = ?", userId, peerId).Delete(&basicModel.Contacts{}).Error
	return err
}

func (d *contactsDao) UnscopedDeleteByGroupId(db *gorm.DB, groupId string) error {
	return db.Unscoped().Where("peer_id = ? AND peer_type = ?", groupId, constant.ContactsPeerType.Group).Delete(&basicModel.Contacts{}).Error
}
