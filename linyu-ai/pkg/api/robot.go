package api

import (
	"context"
	"github.com/gin-gonic/gin"
	aiParam "github.com/linyu-im/linyu-server/linyu-ai/pkg/param"
	"github.com/linyu-im/linyu-server/linyu-ai/pkg/schedule"
	aiService "github.com/linyu-im/linyu-server/linyu-ai/pkg/service"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/response"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/route"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/utils"
)

func init() {
	route.RegisterWhite("POST", "/ai/v1/robot/answers", RobotAnswersHandler)
}

// RobotAnswersHandler 机器人问答
func RobotAnswersHandler(c *gin.Context) {
	param := &aiParam.RobotAnswersParam{}
	if !utils.ShouldBindBodyWithJSONAndValidate(c, param) {
		return
	}

	// 获取短期记忆
	shortMemory := aiService.AiMemoryService.GetShortTermMemory(param.SessionId)
	// 模型输入
	in := map[string]any{
		"question":    param.Question,
		"robotId":     param.RobotId,
		"sessionId":   param.SessionId,
		"shortMemory": shortMemory,
	}
	// 获取机器人处理流
	graph, err := schedule.GetRobotGraph(param.RobotId)
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
	response.Ok(c, ret)
}
