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
	route.Register("POST", "/basic/v1/apply/add/friend", ApplyAddFriendHandler)
	route.Register("POST", "/basic/v1/apply/agree/friend", ApplyAgreeFriendHandler)
	route.Register("POST", "/basic/v1/apply/reject", ApplyRejectHandler)
	route.Register("POST", "/basic/v1/apply/cancel", ApplyCancelHandler)
	route.Register("POST", "/basic/v1/apply/list/friend", ApplyFriendListHandler)
	route.Register("POST", "/basic/v1/apply/list/group", ApplyGroupListHandler)
	route.Register("POST", "/basic/v1/apply/add/group", ApplyAddGroupHandler)
}

// ApplyFriendListHandler 好友申请友列表
func ApplyFriendListHandler(c *gin.Context) {
	currentUserId := c.GetString("userId")
	applyList, err := basicService.ApplyService.ApplyFriendList(currentUserId)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Ok(c, applyList)
}

// ApplyGroupListHandler 群聊申请列表
func ApplyGroupListHandler(c *gin.Context) {
	currentUserId := c.GetString("userId")
	applyList, err := basicService.ApplyService.ApplyGroupList(currentUserId)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Ok(c, applyList)
}

// ApplyAddFriendHandler 申请添加好友
func ApplyAddFriendHandler(c *gin.Context) {
	param := &basicParam.ApplyAddFriendParam{}
	if !utils.ShouldBindBodyWithJSONAndValidate(c, param) {
		return
	}
	currentUserId := c.GetString("userId")
	err := basicService.ApplyService.ApplyAddFriend(currentUserId, param)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Ok(c)
}

// ApplyAddGroupHandler 申请加入群聊
func ApplyAddGroupHandler(c *gin.Context) {
	param := &basicParam.ApplyAddGroupParam{}
	if !utils.ShouldBindBodyWithJSONAndValidate(c, param) {
		return
	}
	currentUserId := c.GetString("userId")
	err := basicService.ApplyService.ApplyAddGroup(currentUserId, param)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Ok(c)
}

// ApplyAgreeFriendHandler 申请同意添加好友
func ApplyAgreeFriendHandler(c *gin.Context) {
	param := &basicParam.ApplyAgreeFriendParam{}
	if !utils.ShouldBindBodyWithJSONAndValidate(c, param) {
		return
	}
	currentUserId := c.GetString("userId")
	err := basicService.ApplyService.ApplyAgreeFriend(currentUserId, param)
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
