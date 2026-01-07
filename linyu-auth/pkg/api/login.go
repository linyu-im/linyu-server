package api

import (
	"github.com/gin-gonic/gin"
	"github.com/linyu-im/linyu-server/linyu-auth/pkg/param"
	authService "github.com/linyu-im/linyu-server/linyu-auth/pkg/service"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/response"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/route"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/utils"
)

func init() {
	route.Register("POST", "/auth/v1/login/pwd", PwdLoginHandler, true)
}

// PwdLoginHandler 密码登录
func PwdLoginHandler(c *gin.Context) {
	pwdLoginParam := &param.PwdLoginParam{}
	if !utils.ShouldBindBodyWithJSONAndValidate(c, pwdLoginParam) {
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
