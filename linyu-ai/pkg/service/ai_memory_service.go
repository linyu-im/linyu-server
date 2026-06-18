package service

import (
	"context"
	"slices"

	"github.com/cloudwego/eino/schema"
	basicService "github.com/linyu-im/linyu-server/linyu-basic-service/pkg/service"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/constant"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/db"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/utils"
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
			msg = schema.UserMessage(m.Content.ToString())
		case constant.MessageFromType.Robot:
			msg = schema.AssistantMessage(m.Content.ToString(), nil)
		default:
			msg = schema.UserMessage(m.Content.ToString())
		}
		result = append(result, msg)
	}
	return result
}

func (s *aiMemoryService) GetLongTermMemory(sessionId string, content string) (*schema.Message, error) {
	embedder, err := GetEmbedder()
	if err != nil {
		return nil, err
	}
	// 获取向量
	v, err := embedder.EmbedStrings(context.Background(), []string{content})
	if err != nil {
		return nil, err
	}
	embeddings := utils.Float64ToFloat32(v[0])
	filter := map[string]string{"sessionId": sessionId}
	results, err := db.Vector.Search(constant.VectorCollection.LongTermMemory, embeddings, filter, 5, 0.3)
	if err != nil {
		return nil, err
	}
	// 长期记忆拼接
	var memory string
	if len(results) > 0 {
		memory = "相关的长期记忆:\n"
		for _, m := range results {
			memory += "-" + m.Content + "\n"
		}
	}
	return schema.SystemMessage(memory), err
}

func (s *aiMemoryService) SaveLongTermMemory(sessionId string, content string) error {
	embedder, err := GetEmbedder()
	if err != nil {
		return err
	}
	v, err := embedder.EmbedStrings(context.Background(), []string{content})
	if err != nil {
		return err
	}
	embeddings := utils.Float64ToFloat32(v[0])
	metadata := map[string]string{"sessionId": sessionId, "content": content}
	err = db.Vector.Insert(constant.VectorCollection.LongTermMemory, embeddings, metadata)
	return err
}
