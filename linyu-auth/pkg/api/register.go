package api

import (
	"github.com/gin-gonic/gin"
	"github.com/linyu-im/linyu-server/linyu-auth/pkg/param"
	basicService "github.com/linyu-im/linyu-server/linyu-basic-service/pkg/service"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/response"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/route"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/utils"
)

func init() {
	route.RegisterWhite("POST", "/auth/v1/register/account/check", AccountCheckHandler)
	route.RegisterWhite("POST", "/auth/v1/register/email/code", RegisterEmailCodeHandler)
	route.RegisterWhite("POST", "/auth/v1/register/email", RegisterByEmailHandler)
}

// AccountCheckHandler 校验账号是否已被使用
func AccountCheckHandler(c *gin.Context) {
	p := &param.AccountCheckParam{}
	if !utils.ShouldBindBodyWithJSONAndValidate(c, p) {
		return
	}
	used := basicService.UserService.GetUserInfoByAccount(p.Account) != nil
	response.Ok(c, &used)
}

// RegisterEmailCodeHandler 注册发送邮箱验证码
func RegisterEmailCodeHandler(c *gin.Context) {
	p := &param.EmailCodeParam{}
	if !utils.ShouldBindBodyWithJSONAndValidate(c, p) {
		return
	}
	user, _ := basicService.UserService.GetUserByKV("email", p.Email)
	if user != nil {
		response.Fail(c, "auth.email-exists")
		return
	}
	if err := basicService.UserService.SendCodeByEmail(p.Email); err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Ok(c)
}

// RegisterByEmailHandler 邮箱方式注册账号
func RegisterByEmailHandler(c *gin.Context) {
	p := &param.EmailRegisterParam{}
	if !utils.ShouldBindBodyWithJSONAndValidate(c, p) {
		return
	}
	user, _ := basicService.UserService.GetUserByKV("email", p.Email)
	if user != nil {
		response.Fail(c, "auth.email-exists")
		return
	}
	if basicService.UserService.GetUserInfoByAccount(p.Account) != nil {
		response.Fail(c, "auth.account-exists")
		return
	}
	if !basicService.UserService.VerifyCode(p.Email, p.Code) {
		response.Fail(c, "auth.code-expire")
		return
	}
	if err := basicService.UserService.RegisterByEmail(p.Email, p.Account, p.Password); err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Ok(c)
}
