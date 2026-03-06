package service

import (
	"context"
	"fmt"
	"github.com/cloudwego/eino-ext/components/model/openai"
	"sync"
)

var ModelPool sync.Map

func GetLLModel(modelId string) (*openai.ChatModel, error) {
	if m, ok := ModelPool.Load(modelId); ok {
		return m.(*openai.ChatModel), nil
	}
	//获取模型信息
	modelInfo := AiModelService.GetAiModel(modelId)
	if modelInfo == nil {
		return nil, fmt.Errorf("model not found")
	}
	ctx := context.Background()
	llm, err := openai.NewChatModel(ctx, &openai.ChatModelConfig{
		BaseURL:             modelInfo.BaseURL,
		Model:               modelInfo.Model,
		APIKey:              modelInfo.APIKey,
		Temperature:         modelInfo.Temperature,
		MaxCompletionTokens: modelInfo.MaxTokens,
	})
	ModelPool.Store(modelId, llm)
	return llm, err
}
