package service

import (
	"errors"
	"fmt"
	basicDao "github.com/linyu-im/linyu-server/linyu-basic-service/internal/dao"
	basicModel "github.com/linyu-im/linyu-server/linyu-basic-service/pkg/model"
	basicParam "github.com/linyu-im/linyu-server/linyu-basic-service/pkg/param"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/constant"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/db"
	emailutil "github.com/linyu-im/linyu-server/linyu-common/pkg/email"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/localtime"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/request"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/response"
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

func (s *userService) CreateUserByKV(kvs ...interface{}) (*basicModel.User, error) {
	if len(kvs)%2 != 0 {
		return nil, fmt.Errorf("invalid number of key-value pairs: must be even, got %d", len(kvs))
	}

	// 生成唯一账号
	account := utils.GenerateOnlyNumber("linyu_", func(account string) bool {
		user := basicDao.UserDao.GetUserByAccount(db.RDB, account)
		return user == nil
	})

	//构建基础用户信息map
	userMap := map[string]interface{}{
		"id":         utils.GenerateSfIDString(),
		"username":   utils.RandUsername("林语"),
		"account":    account,
		"created_at": localtime.Now(),
		"updated_at": localtime.Now(),
	}

	// 遍历传入的所有key-value并添加到map中
	for i := 0; i < len(kvs); i += 2 {
		key, ok := kvs[i].(string)
		if !ok {
			return nil, fmt.Errorf("key at position %d is not a string", i)
		}
		value := kvs[i+1]
		userMap[key] = value
	}

	//创建用户
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

func (s *userService) GetUserByGitee(gitee string) *basicModel.User {
	user := basicDao.UserDao.GetUserByGitee(db.RDB, gitee)
	return user
}

func (s *userService) GetUserById(userId string) *basicModel.User {
	user := basicDao.UserDao.GetUserById(db.RDB, userId)
	return user
}

func (s *userService) CurrentUserInfoById(userId string) (*basicModel.User, error) {
	return basicDao.UserDao.CurrentUserInfoById(db.RDB, userId)
}

func (s *userService) UserInfoById(userId string, currentUserId string) (*basicModel.User, error) {
	user, err := basicDao.UserDao.UserInfoById(db.RDB, userId, currentUserId)
	if err != nil {
		return nil, err
	}
	tx := basicDao.MomentDao.BuildMomentQuery(db.RDB, userId, userId)
	pages, err := response.Paginate[*basicModel.Moment](tx, request.PageQuery{Page: 1, PageSize: 1})
	if err == nil && len(pages.Records) > 0 {
		user.Moment = pages.Records[0]
	}
	return user, nil
}

func (s *userService) GetAvatar(userId string) string {
	return basicDao.UserDao.UserAvatarById(db.RDB, userId)
}

func (s *userService) UpdateAvatar(userId string, avatarUrl string) error {
	return basicDao.UserDao.UpdateAvatar(db.RDB, userId, avatarUrl)
}

func (s *userService) UpdateProfile(userId string, param *basicParam.UserUpdateProfileParam) error {
	fields := map[string]interface{}{}
	if param.Username != "" {
		fields["username"] = param.Username
	}
	if param.Gender != "" {
		fields["gender"] = param.Gender
	}
	if param.Birthday != "" {
		fields["birthday"] = param.Birthday
	}
	if param.Signature != "" {
		fields["signature"] = param.Signature
	}
	if param.Location != "" {
		fields["location"] = param.Location
	}
	if len(fields) == 0 {
		return nil
	}
	return basicDao.UserDao.UpdateProfile(db.RDB, userId, fields)
}
