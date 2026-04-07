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
