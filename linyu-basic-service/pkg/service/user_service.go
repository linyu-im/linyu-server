package service

import (
	"errors"
	"fmt"
	basicDao "github.com/linyu-im/linyu-server/linyu-basic-service/internal/dao"
	basicModel "github.com/linyu-im/linyu-server/linyu-basic-service/pkg/model"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/constant"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/db"
	emailutil "github.com/linyu-im/linyu-server/linyu-common/pkg/email"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/localtime"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/utils"
	"time"
)

var UserService = newUserService()

func newUserService() *userService {
	return &userService{}
}

type userService struct{}

// GetUserInfoByAccount 根据账号获取用户信息
func (s *userService) GetUserInfoByAccount(account string) *basicModel.User {
	user := basicDao.UserDao.GetUserByAccount(db.RDB, account)
	return user
}

// VerifyCode 校验验证码
func (s *userService) VerifyCode(tag string, code string) bool {
	key := fmt.Sprintf(constant.RedisKey.UserCode, tag)
	codeRedis, err := db.CacheDB.Get(key)
	if err != nil || code != codeRedis {
		return false
	}
	_ = db.CacheDB.Del(key)
	return true
}

func (s *userService) SendCodeByEmail(email string) error {
	code, err := s.GenerateCode(email)
	if err != nil {
		return err
	}
	emailutil.SendEmailCode(email, code)
	return nil
}

// RegisterByEmail 根据邮箱创建账号
func (s *userService) RegisterByEmail(email string) error {
	account := utils.GenerateOnlyNumber("linyu_", func(account string) bool {
		user := basicDao.UserDao.GetUserByAccount(db.RDB, account)
		return user == nil
	})
	user := &basicModel.User{
		ID:       utils.GenerateSfIDString(),
		Email:    &email,
		Username: utils.RandUsername("林语"),
		Account:  account,
	}
	err := basicDao.UserDao.Create(db.RDB, user)
	if err != nil {
		return err
	}
	return nil
}

func (s *userService) CreateUserByKV(key, value string) (*basicModel.User, error) {
	account := utils.GenerateOnlyNumber("linyu_", func(account string) bool {
		user := basicDao.UserDao.GetUserByAccount(db.RDB, account)
		return user == nil
	})
	userMap := map[string]interface{}{
		"id":         utils.GenerateSfIDString(),
		"username":   utils.RandUsername("林语"),
		"account":    account,
		"created_at": localtime.Now(),
		"updated_at": localtime.Now(),
		key:          value,
	}
	err := basicDao.UserDao.CreateByMap(db.RDB, userMap)
	if err != nil {
		return nil, err
	}
	return basicDao.UserDao.GetUserByAccount(db.RDB, account), nil
}

// GenerateCode 生成验证码
func (s *userService) GenerateCode(tag string) (string, error) {
	//60s内，只能发送一次
	lock, err := db.CacheDB.Exists(fmt.Sprintf(constant.RedisKey.UserCodeLock, tag))
	if err != nil {
		return "", err
	}
	if lock {
		return "", errors.New("auth.code-send-too-frequent")
	}
	code := utils.Random6DigitCode()
	//验证码10分钟内有效
	if err := db.CacheDB.Set(fmt.Sprintf(constant.RedisKey.UserCode, tag), code, 10*time.Minute); err != nil {
		return "", err
	}
	if err := db.CacheDB.Set(fmt.Sprintf(constant.RedisKey.UserCodeLock, tag), 1, 60*time.Second); err != nil {
		return "", err
	}
	return code, nil
}

func (s *userService) GetUserByKV(key string, value string) (*basicModel.User, error) {
	return basicDao.UserDao.GetUserByKV(db.RDB, key, value), nil
}
