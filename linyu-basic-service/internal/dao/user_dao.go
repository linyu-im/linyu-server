package dao

import (
	"fmt"

	basicModel "github.com/linyu-im/linyu-server/linyu-basic-service/pkg/model"
	"gorm.io/gorm"
)

var UserDao = newUserDao()

func newUserDao() *userDao {
	return &userDao{}
}

type userDao struct{}

// GetUserByAccount 根据账号获取用户信息
func (r userDao) GetUserByAccount(db *gorm.DB, account string) *basicModel.User {
	result := &basicModel.User{}
	if err := db.First(result, "account = ?", account).Error; err != nil {
		return nil
	}
	return result
}

func (r userDao) GetUserByEmail(db *gorm.DB, email string) *basicModel.User {
	result := &basicModel.User{}
	if err := db.First(result, "email = ?", email).Error; err != nil {
		return nil
	}
	return result
}

func (r userDao) Create(db *gorm.DB, user *basicModel.User) error {
	if err := db.Create(user).Error; err != nil {
		return err
	}
	return nil
}

func (r userDao) CreateByMap(db *gorm.DB, userMap map[string]interface{}) error {
	if err := db.Model(&basicModel.User{}).Create(userMap).Error; err != nil {
		return err
	}
	return nil
}

func (r userDao) GetUserByKV(db *gorm.DB, key string, value string) *basicModel.User {
	result := &basicModel.User{}
	err := db.Where(fmt.Sprintf("%s = ?", key), value).
		First(result).Error
	if err != nil {
		return nil
	}
	return result
}

func (r userDao) GetUserByGitee(db *gorm.DB, gitee string) *basicModel.User {
	result := &basicModel.User{}
	if err := db.First(result, "gitee = ?", gitee).Error; err != nil {
		return nil
	}
	return result
}

func (r userDao) GetUserById(db *gorm.DB, userId string) *basicModel.User {
	result := &basicModel.User{}
	if err := db.First(result, "id = ?", userId).Error; err != nil {
		return nil
	}
	return result
}

func (r userDao) CurrentUserInfoById(db *gorm.DB, userId string) (*basicModel.User, error) {
	var user basicModel.User
	err := db.
		Table("t_user u").
		Select("u.*, e.id as emotion_id, e.emotion_name, e.url as emotion_url, e.type").
		Joins("LEFT JOIN t_emotion e ON u.emotion_id = e.id").
		Where("u.id = ?", userId).
		Scan(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r userDao) UserAvatarById(rdb *gorm.DB, userId string) string {
	var avatarUrl string
	err := rdb.
		Table("t_user u").
		Select("avatar").
		Where("id = ?", userId).
		Scan(&avatarUrl).Error
	if err != nil {
		return ""
	}
	return avatarUrl
}

func (r userDao) UpdatePassword(db *gorm.DB, userId string, password string) error {
	return db.Model(&basicModel.User{}).
		Where("id = ?", userId).
		Update("password", password).Error
}

func (r userDao) UpdateAvatar(db *gorm.DB, userId string, avatar string) error {
	return db.Model(&basicModel.User{}).
		Where("id = ?", userId).
		Update("avatar", avatar).Error
}

func (r userDao) UpdateProfile(db *gorm.DB, userId string, fields map[string]interface{}) error {
	return db.Model(&basicModel.User{}).
		Where("id = ?", userId).
		Updates(fields).Error
}

func (r userDao) SearchByKeyword(db *gorm.DB, keyword string) *gorm.DB {
	like := "%" + keyword + "%"
	return db.Model(&basicModel.User{}).Where("account LIKE ?", like)
}

func (r userDao) UserInfoById(db *gorm.DB, userId string, currentUserId string) (*basicModel.User, error) {
	var user basicModel.User
	err := db.
		Table("t_user u").
		Select("u.*, e.id as emotion_id, e.emotion_name, e.url as emotion_url, e.type, c.remark, c.tag").
		Joins("LEFT JOIN t_emotion e ON u.emotion_id = e.id").
		Joins("LEFT JOIN t_contacts c ON c.user_id = ? AND c.peer_id = u.id", currentUserId).
		Where("u.id = ?", userId).
		Scan(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
