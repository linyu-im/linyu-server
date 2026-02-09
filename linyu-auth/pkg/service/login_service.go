package service

import (
	"errors"
	"fmt"
	"github.com/go-ldap/ldap/v3"
	authResult "github.com/linyu-im/linyu-server/linyu-auth/pkg/result"
	basicModel "github.com/linyu-im/linyu-server/linyu-basic-service/pkg/model"
	basicService "github.com/linyu-im/linyu-server/linyu-basic-service/pkg/service"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/config"
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
	err = db.CacheDB.Set(fmt.Sprintf(constant.RedisKey.UserToken, user.ID, device), loginVersion, jwt.GetJwtExpireTime())
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

func (s *loginService) LdapLogin(username, pwd, device string) (*authResult.UserLoginInfoResult, error) {
	// 获取ldap验证内容
	uniqueValue, err := LdapAuthenticate(username, pwd)
	if err != nil {
		return nil, err
	}
	user, err := basicService.UserService.GetUserByKV(config.C.Auth.Ldap.Unique.LocalField, uniqueValue)
	if user == nil {
		//创建用户
		user, err = basicService.UserService.CreateUserByKV(config.C.Auth.Ldap.Unique.LocalField, uniqueValue)
		if err != nil {
			return nil, err
		}
	}
	return s.Login(user, device)
}

// LdapAuthenticate ldap验证
func LdapAuthenticate(username, password string) (string, error) {
	if !config.C.Auth.Ldap.Enabled {
		return "", errors.New("auth.ldap-not-enabled")
	}
	conn, err := ldap.DialURL(config.C.Auth.Ldap.Host)
	if err != nil {
		return "", err
	}
	defer conn.Close()

	if err = conn.Bind(config.C.Auth.Ldap.BindDN, config.C.Auth.Ldap.BindPassword); err != nil {
		return "", err
	}
	filter := fmt.Sprintf(config.C.Auth.Ldap.UserFilter, username)

	searchRequest := ldap.NewSearchRequest(
		config.C.Auth.Ldap.BaseDN,
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		0, 0, false,
		filter,
		nil,
		nil,
	)
	sr, err := conn.Search(searchRequest)
	if err != nil || len(sr.Entries) == 0 {
		return "", errors.New("auth.user-not-exist")
	}
	entry := sr.Entries[0]
	if err = conn.Bind(entry.DN, password); err != nil {
		return "", errors.New("auth.password-error")
	}
	return entry.GetAttributeValue(config.C.Auth.Ldap.Unique.LdapField), nil
}
