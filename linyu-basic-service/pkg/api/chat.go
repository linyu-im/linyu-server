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
	route.Register("POST", "/basic/v1/chat/top", ChatSetTopHandler)
	route.Register("POST", "/basic/v1/chat/mute", ChatMuteHandler)
	route.Register("POST", "/basic/v1/chat/delete", ChatDeleteHandler)
	route.Register("POST", "/basic/v1/chat/mark-read", ChatMarkReadHandler)
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

// ChatMuteHandler 聊天会话免打扰设置
func ChatMuteHandler(c *gin.Context) {
	param := &basicParam.ChatMuteParam{}
	if !utils.ShouldBindBodyWithJSONAndValidate(c, param) {
		return
	}
	currentUserId := c.GetString("userId")
	err := basicService.ChatService.SetMute(currentUserId, param)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Ok(c)
}

// ChatDeleteHandler 聊天会话删除
func ChatDeleteHandler(c *gin.Context) {
	param := &basicParam.ChatDeleteParam{}
	if !utils.ShouldBindBodyWithJSONAndValidate(c, param) {
		return
	}
	currentUserId := c.GetString("userId")
	err := basicService.ChatService.ChatDelete(currentUserId, param)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Ok(c)
}

// ChatMarkReadHandler 聊天会话已读
func ChatMarkReadHandler(c *gin.Context) {
	param := &basicParam.ChatMarkReadParam{}
	if !utils.ShouldBindBodyWithJSONAndValidate(c, param) {
		return
	}
	currentUserId := c.GetString("userId")
	err := basicService.ChatService.MarkRead(currentUserId, param)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Ok(c)
}
