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
	route.RegisterWhite("POST", "/basic/v1/app/version/check", AppVersionCheckHandler)
}

// AppVersionCheckHandler 客户端版本检测
func AppVersionCheckHandler(c *gin.Context) {
	param := &basicParam.AppVersionCheckParam{}
	if !utils.ShouldBindBodyWithJSONAndValidate(c, param) {
		return
	}
	result, err := basicService.AppVersionService.Check(param)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Ok(c, result)
}
