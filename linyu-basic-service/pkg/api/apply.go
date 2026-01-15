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
	route.Register("POST", "/basic/v1/apply/add/contacts", ApplyAddContactsHandler, false)
	route.Register("POST", "/basic/v1/apply/agree/contacts", ApplyAgreeContactsHandler, false)
	route.Register("POST", "/basic/v1/apply/reject", ApplyRejectHandler, false)
	route.Register("POST", "/basic/v1/apply/cancel", ApplyCancelHandler, false)
	route.Register("POST", "/basic/v1/apply/list", ApplyListHandler, false)
}

// ApplyListHandler 申请列表
func ApplyListHandler(c *gin.Context) {
	currentUserId := c.GetString("userId")
	applyList, err := basicService.ApplyService.ApplyList(currentUserId)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Ok(c, applyList)
}

// ApplyAddContactsHandler 申请添加联系人
func ApplyAddContactsHandler(c *gin.Context) {
	param := &basicParam.ApplyAddContactsParam{}
	if !utils.ShouldBindBodyWithJSONAndValidate(c, param) {
		return
	}
	currentUserId := c.GetString("userId")
	err := basicService.ApplyService.ApplyAddContacts(currentUserId, param)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Ok(c)
}

// ApplyAgreeContactsHandler 申请同意添加联系人
func ApplyAgreeContactsHandler(c *gin.Context) {
	param := &basicParam.ApplyAgreeContactsParam{}
	if !utils.ShouldBindBodyWithJSONAndValidate(c, param) {
		return
	}
	currentUserId := c.GetString("userId")
	err := basicService.ApplyService.ApplyAgreeContacts(currentUserId, param)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Ok(c)
}

// ApplyRejectHandler 申请拒绝
func ApplyRejectHandler(c *gin.Context) {
	param := &basicParam.ApplyRejectParam{}
	if !utils.ShouldBindBodyWithJSONAndValidate(c, param) {
		return
	}
	currentUserId := c.GetString("userId")
	err := basicService.ApplyService.ApplyReject(currentUserId, param)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Ok(c)
}

// ApplyCancelHandler 申请取消
func ApplyCancelHandler(c *gin.Context) {
	param := &basicParam.ApplyCancelParam{}
	if !utils.ShouldBindBodyWithJSONAndValidate(c, param) {
		return
	}
	currentUserId := c.GetString("userId")
	err := basicService.ApplyService.ApplyCancel(currentUserId, param)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Ok(c)
}
