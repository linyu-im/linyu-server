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
	route.Register("POST", "/basic/v1/message/send/user", SendMessageToUserHandler)
	route.Register("POST", "/basic/v1/message/send/group", SendMessageToGroupHandler)
	route.Register("POST", "/basic/v1/message/page", MessagePageHandler)
}

// SendMessageToUserHandler 发送消息(用户)
func SendMessageToUserHandler(c *gin.Context) {
	param := &basicParam.SendMessageToUserParam{}
	if !utils.ShouldBindBodyWithJSONAndValidate(c, param) {
		return
	}
	currentUserId := c.GetString("userId")
	msg, err := basicService.MessageService.SendMessageToUser(currentUserId, param)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Ok(c, msg)
}

// SendMessageToGroupHandler 发送消息(群聊)
func SendMessageToGroupHandler(c *gin.Context) {
	param := &basicParam.SendMessageToGroupParam{}
	if !utils.ShouldBindBodyWithJSONAndValidate(c, param) {
		return
	}
	currentUserId := c.GetString("userId")
	msg, err := basicService.MessageService.SendMessageToGroup(currentUserId, param)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Ok(c, msg)
}

// MessagePageHandler 分页获取聊天内容
func MessagePageHandler(c *gin.Context) {
	param := &basicParam.MessagePageParam{}
	if !utils.ShouldBindBodyWithJSONAndValidate(c, param) {
		return
	}
	currentUserId := c.GetString("userId")
	data, err := basicService.MessageService.MessagePage(currentUserId, param)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Ok(c, data)
}
