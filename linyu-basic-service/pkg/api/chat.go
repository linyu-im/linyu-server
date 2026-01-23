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
	route.Register("POST", "/basic/v1/chat/list", ChatListHandler)
	route.Register("POST", "/basic/v1/chat/create", ChatCreateHandler)
	route.Register("POST", "/basic/v1/chat/set/top", ChatSetTopHandler)
}

// ChatListHandler 聊天会话列表
func ChatListHandler(c *gin.Context) {
	currentUserId := c.GetString("userId")
	list, err := basicService.ChatService.ChatList(currentUserId)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Ok(c, list)
}

// ChatCreateHandler 聊天会话创建
func ChatCreateHandler(c *gin.Context) {
	param := &basicParam.ChatCreateParam{}
	if !utils.ShouldBindBodyWithJSONAndValidate(c, param) {
		return
	}
	currentUserId := c.GetString("userId")
	chat, err := basicService.ChatService.ChatCreate(currentUserId, param)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Ok(c, chat)
}

// ChatSetTopHandler 聊天会话置顶设置
func ChatSetTopHandler(c *gin.Context) {
	param := &basicParam.ChatSetTopParam{}
	if !utils.ShouldBindBodyWithJSONAndValidate(c, param) {
		return
	}
	currentUserId := c.GetString("userId")
	err := basicService.ChatService.SetTop(currentUserId, param)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Ok(c)
}
