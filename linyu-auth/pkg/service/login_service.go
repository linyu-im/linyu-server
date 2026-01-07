package service

import (
	"errors"
	"fmt"
	authResult "github.com/linyu-im/linyu-server/linyu-auth/pkg/result"
	basicModel "github.com/linyu-im/linyu-server/linyu-basic-service/pkg/model"
	basicService "github.com/linyu-im/linyu-server/linyu-basic-service/pkg/service"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/constant"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/db"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/jwt"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/utils"
)

var LoginService = newLoginService()

func newLoginService() *loginService {
	return &loginService{}
}

type loginService struct{}

func (s *loginService) Login(user *basicModel.User, device string) (*authResult.UserLoginInfoResult, error) {
	loginVersion := utils.GenerateUuid()
	userInfo := jwt.JwtClaims{
		UserID:       user.ID,
		LoginVersion: loginVersion,
	}
	if user.Status == constant.UserStatus.Banned {
		return nil, errors.New("auth.user-banned")
	}
	token, err := jwt.GenerateJwtToken(userInfo)
	if err != nil {
		return nil, errors.New("auth.error")
	}
	err = db.RedisDB.Set(fmt.Sprintf(constant.RedisKey.UserToken, user.ID, device), loginVersion, jwt.GetJwtExpireTime())
	if err != nil {
		return nil, errors.New("auth.error")
	}
	result := &authResult.UserLoginInfoResult{
		UserID: user.ID,
		Token:  token,
	}
	return result, nil
}

func (s *loginService) PasswordLogin(account, pwd, device string) (*authResult.UserLoginInfoResult, error) {
	user := basicService.UserService.GetUserInfoByAccount(account)
	if user == nil {
		return nil, errors.New("auth.user-not-exist")
	}
	if b, _ := utils.VerifyPasswordArgon2id(pwd, user.Password); !b {
		return nil, errors.New("auth.password-error")
	}
	return s.Login(user, device)
}
