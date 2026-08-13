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
	route.RegisterWhite("POST", "/auth/v1/password/email/code", ResetPasswordEmailCodeHandler)
	route.RegisterWhite("POST", "/auth/v1/password/reset", ResetPasswordHandler)
}

// ResetPasswordEmailCodeHandler 找回密码发送邮箱验证码
func ResetPasswordEmailCodeHandler(c *gin.Context) {
	p := &param.ResetPasswordEmailCodeParam{}
	if !utils.ShouldBindBodyWithJSONAndValidate(c, p) {
		return
	}
	if err := basicService.UserService.SendResetPasswordCodeByEmail(p.Email); err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Ok(c)
}

// ResetPasswordHandler 邮箱验证码找回/重置密码
func ResetPasswordHandler(c *gin.Context) {
	p := &param.ResetPasswordParam{}
	if !utils.ShouldBindBodyWithJSONAndValidate(c, p) {
		return
	}
	if err := basicService.UserService.ResetPasswordByEmail(p.Email, p.Code, p.Password); err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Ok(c)
}
