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
	route.Register("POST", "/basic/v1/application/list", ApplicationListHandler)
}

// ApplicationListHandler 查询应用列表
func ApplicationListHandler(c *gin.Context) {
	param := &basicParam.ApplicationListParam{}
	if !utils.ShouldBindBodyWithJSONAndValidate(c, param) {
		return
	}
	list, err := basicService.ApplicationService.ApplicationList(param)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Ok(c, list)
}
