package api

import (
	"github.com/gin-gonic/gin"
	appParam "github.com/linyu-im/linyu-server/linyu-application/pkg/param"
	appService "github.com/linyu-im/linyu-server/linyu-application/pkg/service"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/response"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/route"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/utils"
)

func init() {
	route.Register("POST", "/application/v1/list", ApplicationListHandler)
}

// ApplicationListHandler 查询应用列表
func ApplicationListHandler(c *gin.Context) {
	param := &appParam.ApplicationListParam{}
	if !utils.ShouldBindBodyWithJSONAndValidate(c, param) {
		return
	}
	list, err := appService.ApplicationService.ApplicationList(param)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Ok(c, list)
}
