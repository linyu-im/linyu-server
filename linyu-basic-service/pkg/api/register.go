package api

import (
	"github.com/gin-gonic/gin"
	basicParam "github.com/linyu-im/linyu-server/linyu-basic-service/pkg/param"
	basicService "github.com/linyu-im/linyu-server/linyu-basic-service/pkg/service"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/response"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/route"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/utils"
)

func init() {
	route.Register("POST", "/basic/v1/register/email", EmailRegisterHandler, true)
	route.Register("POST", "/basic/v1/code/email", SendEmailCodeHandler, true)
}

// SendEmailCodeHandler 发送邮件验证码
func SendEmailCodeHandler(c *gin.Context) {
	param := &basicParam.EmailCodeParam{}
	if !utils.ShouldBindBodyWithJSONAndValidate(c, param) {
		return
	}
	err := basicService.UserService.SendCodeByEmail(param.Email)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Ok(c)
}

// EmailRegisterHandler 邮箱方式注册账号
func EmailRegisterHandler(c *gin.Context) {
	emailRegParam := &basicParam.EmailRegisterParam{}
	if !utils.ShouldBindBodyWithJSONAndValidate(c, emailRegParam) {
		return
	}
	//校验验证码是否有效
	result := basicService.UserService.VerifyCode(emailRegParam.Email, emailRegParam.Code)
	if !result {
		response.Fail(c, "auth.code-expire")
		return
	}
	//注册用户
	err := basicService.UserService.RegisterByEmail(emailRegParam.Email)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Ok(c)
}
