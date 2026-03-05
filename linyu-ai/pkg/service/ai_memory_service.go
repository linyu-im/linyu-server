package service

import (
	"github.com/cloudwego/eino/schema"
	basicService "github.com/linyu-im/linyu-server/linyu-basic-service/pkg/service"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/constant"
	"slices"
)

var AiMemoryService = newAiMemoryService()

func newAiMemoryService() *aiMemoryService {
	return &aiMemoryService{}
}

type aiMemoryService struct{}

func (s *aiMemoryService) GetShortTermMemory(sessionId string) []*schema.Message {
	if len(sessionId) == 0 {
		return []*schema.Message{}
	}
	messages := basicService.MessageService.GetMessageBySessionId(sessionId, 10)
	//翻转消息数组
	slices.Reverse(messages)
	result := make([]*schema.Message, 0, len(messages))
	for _, m := range messages {
		var msg *schema.Message
		switch m.FromType {
		case constant.MessageFromType.User:
			msg = schema.UserMessage(m.Content)
		case constant.MessageFromType.Bot:
			msg = schema.AssistantMessage(m.Content, nil)
		default:
			msg = schema.UserMessage(m.Content)
		}
		result = append(result, msg)
	}
	return result
}
