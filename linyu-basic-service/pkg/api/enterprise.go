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
	route.Register("POST", "/basic/v1/enterprise/info", EnterpriseInfoHandler)
	route.Register("POST", "/basic/v1/enterprise/avatar/get", GetEnterpriseAvatarHandler)
	route.Register("POST", "/basic/v1/enterprise/department", EnterpriseDepartmentHandler)
}

// EnterpriseInfoHandler 获取企业详情
func EnterpriseInfoHandler(c *gin.Context) {
	param := &basicParam.EnterpriseInfoParam{}
	if !utils.ShouldBindBodyWithJSONAndValidate(c, param) {
		return
	}
	currentUserId := c.GetString("userId")
	enterprise, err := basicService.EnterpriseService.EnterpriseInfo(currentUserId, param.EnterpriseID)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Ok(c, enterprise)
}

func GetEnterpriseAvatarHandler(c *gin.Context) {
	param := &basicParam.GetEnterpriseAvatarParam{}
	if !utils.ShouldBindBodyWithJSONAndValidate(c, param) {
		return
	}
	url := basicService.EnterpriseService.GetEnterpriseAvatar(param.EnterpriseID)
	response.Ok(c, url)
}

// EnterpriseDepartmentHandler 获取企业部门树（含子部门及部门成员）
func EnterpriseDepartmentHandler(c *gin.Context) {
	param := &basicParam.EnterpriseInfoParam{}
	if !utils.ShouldBindBodyWithJSONAndValidate(c, param) {
		return
	}
	currentUserId := c.GetString("userId")
	list, err := basicService.EnterpriseService.EnterpriseDepartment(param.EnterpriseID, currentUserId)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Ok(c, list)
}
