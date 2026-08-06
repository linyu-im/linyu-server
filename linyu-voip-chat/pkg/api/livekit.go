package api

import (
	"github.com/gin-gonic/gin"
	basicService "github.com/linyu-im/linyu-server/linyu-basic-service/pkg/service"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/config"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/constant"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/response"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/route"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/utils"
	voipParam "github.com/linyu-im/linyu-server/linyu-voip-chat/pkg/param"
	voipService "github.com/linyu-im/linyu-server/linyu-voip-chat/pkg/service"
)

func init() {
	route.Register("POST", "/voip/v1/livekit/host", LivekitHostHandler)
	route.Register("POST", "/voip/v1/livekit/token/group", LivekitTokenHandler)
	route.Register("POST", "/voip/v1/livekit/token/user", LivekitUserTokenHandler)
	route.Register("POST", "/voip/v1/livekit/room/users", LivekitGroupRoomUsersHandler)
}

func LivekitHostHandler(c *gin.Context) {
	response.Ok(c, config.C.Livekit.Host)
}

func LivekitTokenHandler(c *gin.Context) {
	param := &voipParam.LivekitTokenParam{}
	if !utils.ShouldBindBodyWithJSONAndValidate(c, param) {
		return
	}
	currentUserId := c.GetString("userId")
	// 解析群 id，校验 sessionId 为群聊格式
	groupId := utils.GetGroupIdFromSessionID(param.SessionId)
	if groupId == "" {
		response.Fail(c, "param.error")
		return
	}
	//验证用户是否属于该群聊
	is := basicService.GroupService.IsGroupMember(groupId, currentUserId)
	if !is {
		response.Fail(c, "param.error")
		return
	}
	token, err := voipService.LivekitService.LivekitToken(param.SessionId, currentUserId)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Ok(c, token)
}

// LivekitUserTokenHandler 单聊入会 token，roomId 为单聊 sessionId
func LivekitUserTokenHandler(c *gin.Context) {
	param := &voipParam.LivekitUserTokenParam{}
	if !utils.ShouldBindBodyWithJSONAndValidate(c, param) {
		return
	}
	currentUserId := c.GetString("userId")
	// 解析对端用户 id，校验 sessionId 为当前用户参与的单聊会话
	peerId := utils.GetPeerIdFromUserSession(param.SessionId, currentUserId)
	if peerId == "" {
		response.Fail(c, "param.error")
		return
	}
	// 校验对端用互相是否是好友
	if !basicService.ContactsService.IsFriendBothAnd(currentUserId, peerId) {
		response.Fail(c, "param.error")
		return
	}
	token, err := voipService.LivekitService.LivekitToken(param.SessionId, currentUserId)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Ok(c, token)
}

// LivekitGroupRoomUsersHandler 查询房间在线用户列表
func LivekitGroupRoomUsersHandler(c *gin.Context) {
	param := &voipParam.LivekitTokenParam{}
	if !utils.ShouldBindBodyWithJSONAndValidate(c, param) {
		return
	}
	currentUserId := c.GetString("userId")
	switch utils.GetSessionSceneType(param.SessionId) {
	case constant.SceneType.Group:
		groupId := utils.GetGroupIdFromSessionID(param.SessionId)
		if groupId == "" || !basicService.GroupService.IsGroupMember(groupId, currentUserId) {
			response.Fail(c, "param.error")
			return
		}
	case constant.SceneType.User:
		peerId := utils.GetPeerIdFromUserSession(param.SessionId, currentUserId)
		if peerId == "" || !basicService.ContactsService.IsFriendBothAnd(currentUserId, peerId) {
			response.Fail(c, "param.error")
			return
		}
	default:
		response.Fail(c, "param.error")
		return
	}
	list, err := voipService.LivekitService.ListRoomUsers(param.SessionId)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Ok(c, list)
}
