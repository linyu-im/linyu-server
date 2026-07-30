package api

import (
	"github.com/gin-gonic/gin"
	basicParam "github.com/linyu-im/linyu-server/linyu-basic-service/pkg/param"
	basicService "github.com/linyu-im/linyu-server/linyu-basic-service/pkg/service"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/i18n"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/response"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/route"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/storage"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/utils"
)

func init() {
	route.Register("POST", "/basic/v1/message/send", SendMessageHandler)
	route.Register("POST", "/basic/v1/message/page", MessagePageHandler)
	route.Register("POST", "/basic/v1/message/list", MessageListHandler)
	route.Register("POST", "/basic/v1/message/file/upload", UploadMsgFileChunkHandler)
	route.Register("POST", "/basic/v1/message/file/merge", MergeMsgFileChunkHandler)
	route.Register("POST", "/basic/v1/message/forward", ForwardMessageHandler)
}

// SendMessageHandler 发送消息(会话)
func SendMessageHandler(c *gin.Context) {
	param := &basicParam.SendMessageToUserParam{}
	if !utils.ShouldBindBodyWithJSONAndValidate(c, param) {
		return
	}
	currentUserId := c.GetString("userId")
	msg, err := basicService.MessageService.SendMessage(currentUserId, param)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	if msg.Status == "failed" && msg.FailReason != "" {
		msg.FailReason = i18n.T(c, msg.FailReason, nil)
	}
	response.Ok(c, msg)
}

// ForwardMessageHandler 转发消息
func ForwardMessageHandler(c *gin.Context) {
	param := &basicParam.ForwardMessageParam{}
	if !utils.ShouldBindBodyWithJSONAndValidate(c, param) {
		return
	}
	currentUserId := c.GetString("userId")
	err := basicService.MessageService.ForwardMessage(currentUserId, param)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Ok(c)
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

// MessageListHandler 获取指定消息之后的所有消息
func MessageListHandler(c *gin.Context) {
	param := &basicParam.MessageListParam{}
	if !utils.ShouldBindBodyWithJSONAndValidate(c, param) {
		return
	}
	currentUserId := c.GetString("userId")
	data, err := basicService.MessageService.MessageList(currentUserId, param)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Ok(c, data)
}

// UploadMsgFileChunkHandler 聊天文件切片上传
func UploadMsgFileChunkHandler(c *gin.Context) {
	fileHash := c.PostForm("fileHash")
	chunkIndex := c.PostForm("chunkIndex")
	fileName := c.PostForm("fileName")
	if fileHash == "" || chunkIndex == "" {
		response.Fail(c, "param.error")
		return
	}
	file, err := c.FormFile("file")
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	currentUserId := c.GetString("userId")
	if err = storage.UploadChunk(file, chunkIndex, fileHash, fileName, "msgfile/"+currentUserId); err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Ok(c)
}

// MergeMsgFileChunkHandler 聊天文件切片合并
func MergeMsgFileChunkHandler(c *gin.Context) {
	param := &basicParam.UploadMsgFileInfoParam{}
	if !utils.ShouldBindBodyWithJSONAndValidate(c, param) {
		return
	}
	currentUserId := c.GetString("userId")
	err, storagePath := storage.MergeChunk(param.FileHash, param.TotalChunk, param.FileName, "msgfile/"+currentUserId)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Ok(c, storagePath)
}
