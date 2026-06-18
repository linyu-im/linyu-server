package api

import (
	"context"

	"github.com/gin-gonic/gin"
	aiParam "github.com/linyu-im/linyu-server/linyu-ai/pkg/param"
	aiService "github.com/linyu-im/linyu-server/linyu-ai/pkg/service"
	basicModel "github.com/linyu-im/linyu-server/linyu-basic-service/pkg/model"
	basicService "github.com/linyu-im/linyu-server/linyu-basic-service/pkg/service"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/constant"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/response"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/route"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/utils"
)

func init() {
	route.Register("POST", "/ai/v1/robot/answers", RobotAnswersHandler)
	route.Register("POST", "/ai/v1/robot/list", RobotListHandler)
	route.Register("POST", "/ai/v1/robot/avatar/get", RobotAvatarHandler)
	route.Register("POST", "/ai/v1/robot/info", RobotInfoHandler)
}

// RobotInfoHandler 获取机器人详情
func RobotInfoHandler(c *gin.Context) {
	param := &aiParam.GetRobotAvatarParam{}
	if !utils.ShouldBindBodyWithJSONAndValidate(c, param) {
		return
	}
	info := aiService.AiRobotService.GetAiRobot(param.RobotId)
	response.Ok(c, info)
}

// RobotAvatarHandler 获取机器人头像
func RobotAvatarHandler(c *gin.Context) {
	param := &aiParam.GetRobotAvatarParam{}
	if !utils.ShouldBindBodyWithJSONAndValidate(c, param) {
		return
	}
	url := aiService.AiRobotService.GetAvatar(param.RobotId)
	response.Ok(c, url)
}

// RobotListHandler 查询机器人列表
func RobotListHandler(c *gin.Context) {
	list, err := aiService.AiRobotService.List()
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Ok(c, list)
}

// RobotAnswersHandler 机器人问答
func RobotAnswersHandler(c *gin.Context) {
	param := &aiParam.RobotAnswersParam{}
	if !utils.ShouldBindBodyWithJSONAndValidate(c, param) {
		return
	}
	currentUserId := c.GetString("userId")
	var sessionId string
	switch param.MsgScene {
	case constant.MessageScene.User:
		sessionId = utils.Generate1v1SessionID(currentUserId, param.PeerId)
	case constant.MessageScene.Group:
		sessionId = param.PeerId
	}
	// 获取长期记忆
	longMemory, _ := aiService.AiMemoryService.GetLongTermMemory(sessionId, param.Question)
	// 获取短期记忆
	shortMemory := aiService.AiMemoryService.GetShortTermMemory(sessionId)
	// 模型输入
	in := map[string]any{
		"question":    param.Question,
		"robotId":     param.RobotId,
		"sessionId":   sessionId,
		"shortMemory": shortMemory,
		"longMemory":  longMemory,
	}
	// 获取机器人处理流
	graph, err := aiService.GetRobotGraph(param.RobotId)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	ctx := context.Background()
	ret, err := graph.Invoke(ctx, in)

	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	// 回复消息
	message := &basicModel.Message{
		ID:        utils.GenerateSfIDString(),
		SessionID: sessionId,
		FromID:    param.RobotId,
		ToID:      currentUserId,
		MsgScene:  param.MsgScene,
		MsgType:   constant.MessageType.Text,
		FromType:  constant.MessageFromType.Robot,
		Content:   basicModel.TextContent{Text: ret.Content},
	}
	// 发送消息
	msg, err := basicService.MessageService.SendMessageToSession(currentUserId, sessionId, param.MsgScene, message)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Ok(c, msg)
}
