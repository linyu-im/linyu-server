package api

import (
	"github.com/gin-gonic/gin"
	"github.com/linyu-im/linyu-server/linyu-auth/pkg/param"
	authService "github.com/linyu-im/linyu-server/linyu-auth/pkg/service"
	basicService "github.com/linyu-im/linyu-server/linyu-basic-service/pkg/service"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/config"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/response"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/route"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/utils"
)

func init() {
	route.RegisterWhite("POST", "/auth/v1/login/pwd", PwdLoginHandler)
	route.RegisterWhite("POST", "/auth/v1/login/ldap", LdapLoginHandler)
	route.RegisterWhite("POST", "/auth/v1/login/oauth2", Oauth2LoginHandler)
	route.RegisterWhite("POST", "/auth/v1/login/enable/ldap", EnableLdapHandler)
	route.Register("POST", "/auth/v1/login/token/reset", TokenResetHandler)
	route.Register("POST", "/auth/v1/login/password/verify", VerifyPasswordHandler)
}

func ensureLoginVersion(c *gin.Context, platform string, versionCode int) bool {
	if err := basicService.AppVersionService.EnsureMinSupport(platform, versionCode); err != nil {
		response.Fail(c, err.Error())
		return false
	}
	return true
}

// PwdLoginHandler 密码登录
func PwdLoginHandler(c *gin.Context) {
	pwdLoginParam := &param.PwdLoginParam{}
	if !utils.ShouldBindBodyWithJSONAndValidate(c, pwdLoginParam) {
		return
	}
	if !ensureLoginVersion(c, pwdLoginParam.Platform, pwdLoginParam.VersionCode) {
		return
	}
	userInfo, err := authService.LoginService.PasswordLogin(pwdLoginParam.Account, pwdLoginParam.Password,
		c.GetString("device"))
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Ok(c, userInfo)
}

// Oauth2LoginHandler Oauth2登录
func Oauth2LoginHandler(c *gin.Context) {
	loginParam := &param.Oauth2LoginParam{}
	if !utils.ShouldBindBodyWithJSONAndValidate(c, loginParam) {
		return
	}
	if !ensureLoginVersion(c, loginParam.Platform, loginParam.VersionCode) {
		return
	}
	userInfo, err := authService.LoginService.Oauth2Login(loginParam.Type, loginParam.Code,
		c.GetString("device"))
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Ok(c, userInfo)
}

// LdapLoginHandler ldap方式登录
func LdapLoginHandler(c *gin.Context) {
	ldapParam := &param.LdapLoginParam{}
	if !utils.ShouldBindBodyWithJSONAndValidate(c, ldapParam) {
		return
	}
	if !ensureLoginVersion(c, ldapParam.Platform, ldapParam.VersionCode) {
		return
	}
	userInfo, err := authService.LoginService.LdapLogin(ldapParam.Username, ldapParam.Password,
		c.GetString("device"))
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Ok(c, userInfo)
}

// EnableLdapHandler 是否允许ldap登录
func EnableLdapHandler(c *gin.Context) {
	response.Ok(c, config.C.Auth.Ldap.Enabled)
}

// TokenResetHandler token重置
func TokenResetHandler(c *gin.Context) {
	userId := c.GetString("userId")
	userInfo, err := authService.LoginService.TokenReset(userId, c.GetString("device"))
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Ok(c, userInfo)
}

// VerifyPasswordHandler 校验当前登录用户密码是否正确
func VerifyPasswordHandler(c *gin.Context) {
	p := &param.VerifyPasswordParam{}
	if !utils.ShouldBindBodyWithJSONAndValidate(c, p) {
		return
	}
	ok, err := authService.LoginService.VerifyPassword(c.GetString("userId"), p.Password)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Ok(c, &ok)
}
