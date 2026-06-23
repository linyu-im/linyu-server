package api

import (
	"context"
	"encoding/json"
	"io"
	"strings"

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
	route.Register("POST", "/ai/v1/robot/answers/stream", RobotAnswersStreamHandler)
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
	prepared, err := aiService.AiRobotService.PrepareRobotAnswers(currentUserId, param)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	ctx := context.Background()
	ret, err := prepared.Graph.Invoke(ctx, prepared.Input)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	msg, err := sendRobotAnswerMessage(currentUserId, prepared.SessionId, param, ret.Content)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Ok(c, msg)
}

// RobotAnswersStreamHandler 机器人流式问答
func RobotAnswersStreamHandler(c *gin.Context) {
	param := &aiParam.RobotAnswersParam{}
	if !utils.ShouldBindBodyWithJSONAndValidate(c, param) {
		return
	}
	currentUserId := c.GetString("userId")
	prepared, err := aiService.AiRobotService.PrepareRobotAnswers(currentUserId, param)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	ctx := c.Request.Context()
	stream, err := prepared.Graph.Stream(ctx, prepared.Input)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	defer stream.Close()

	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("X-Accel-Buffering", "no")

	var fullContent strings.Builder
	for {
		chunk, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			c.SSEvent("error", gin.H{"msg": err.Error()})
			return
		}
		if chunk.Content == "" {
			continue
		}
		fullContent.WriteString(chunk.Content)
		c.SSEvent("delta", gin.H{"content": chunk.Content})
		c.Writer.Flush()
	}

	msg, err := sendRobotAnswerMessage(currentUserId, prepared.SessionId, param, fullContent.String())
	if err != nil {
		c.SSEvent("error", gin.H{"msg": err.Error()})
		return
	}
	c.SSEvent("done", msg)
}

func sendRobotAnswerMessage(currentUserId, sessionId string, param *aiParam.RobotAnswersParam, content string) (*basicModel.Message, error) {
	contentRaw, err := json.Marshal(basicModel.TextContent{Text: content})
	if err != nil {
		return nil, err
	}
	message := &basicModel.Message{
		ID:        utils.GenerateSfIDString(),
		SessionID: sessionId,
		FromID:    param.RobotId,
		ToID:      currentUserId,
		SceneType: param.SceneType,
		MsgType:   constant.MessageType.Text,
		FromType:  constant.MessageFromType.Robot,
		Content:   contentRaw,
	}
	return basicService.MessageService.SendMessageToSession(currentUserId, message)
}
