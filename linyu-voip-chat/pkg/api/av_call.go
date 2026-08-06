package api

import (
	"github.com/gin-gonic/gin"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/response"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/route"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/utils"
	voipParam "github.com/linyu-im/linyu-server/linyu-voip-chat/pkg/param"
	voipService "github.com/linyu-im/linyu-server/linyu-voip-chat/pkg/service"
)

func init() {
	route.Register("POST", "/voip/v1/av-call/user/invite", AvCallFriendInviteHandler)
	route.Register("POST", "/voip/v1/av-call/user/hangup", AvCallUserHangupHandler)
	route.Register("POST", "/voip/v1/av-call/group/invite", AvCallGroupInviteHandler)
	route.Register("POST", "/voip/v1/av-call/group/hangup", AvCallGroupHangupHandler)
}

// AvCallFriendInviteHandler 好友通话邀请
func AvCallFriendInviteHandler(c *gin.Context) {
	param := &voipParam.AvCallFriendInviteParam{}
	if !utils.ShouldBindBodyWithJSONAndValidate(c, param) {
		return
	}
	currentUserId := c.GetString("userId")
	result, err := voipService.AvCallService.FriendInvite(currentUserId, param)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Ok(c, result)
}

// AvCallUserHangupHandler 单聊通话挂断
func AvCallUserHangupHandler(c *gin.Context) {
	param := &voipParam.AvCallUserHangupParam{}
	if !utils.ShouldBindBodyWithJSONAndValidate(c, param) {
		return
	}
	currentUserId := c.GetString("userId")
	if err := voipService.AvCallService.UserHangup(currentUserId, param); err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Ok(c, nil)
}

// AvCallGroupInviteHandler 群聊通话邀请
func AvCallGroupInviteHandler(c *gin.Context) {
	param := &voipParam.AvCallGroupInviteParam{}
	if !utils.ShouldBindBodyWithJSONAndValidate(c, param) {
		return
	}
	currentUserId := c.GetString("userId")
	result, err := voipService.AvCallService.GroupInvite(currentUserId, param)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Ok(c, result)
}

// AvCallGroupHangupHandler 群聊通话挂断
func AvCallGroupHangupHandler(c *gin.Context) {
	param := &voipParam.AvCallGroupHangupParam{}
	if !utils.ShouldBindBodyWithJSONAndValidate(c, param) {
		return
	}
	currentUserId := c.GetString("userId")
	if err := voipService.AvCallService.GroupHangup(currentUserId, param); err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Ok(c, nil)
}
