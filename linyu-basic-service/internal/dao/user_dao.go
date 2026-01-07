package dao

import (
	"github.com/linyu-im/linyu-server/linyu-basic-service/pkg/model"
	"gorm.io/gorm"
)

var UserDao = newUserDao()

func newUserDao() *userDao {
	return &userDao{}
}

type userDao struct{}

// GetUserByAccount 根据账号获取用户信息
func (r userDao) GetUserByAccount(db *gorm.DB, account string) *model.User {
	result := &model.User{}
	if err := db.First(result, "account = ?", account).Error; err != nil {
		return nil
	}
	return result
}

func (r userDao) GetUserByEmail(db *gorm.DB, email string) *model.User {
	result := &model.User{}
	if err := db.First(result, "email = ?", email).Error; err != nil {
		return nil
	}
	return result
}

func (r userDao) Create(db *gorm.DB, user *model.User) error {
	if err := db.Create(user).Error; err != nil {
		return err
	}
	return nil
}
