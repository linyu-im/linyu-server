package service

import (
	"context"
	"fmt"
	"sync"

	"github.com/cloudwego/eino-ext/components/model/openai"
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
	if err != nil {
		return nil, err
	}
	ModelPool.Store(modelId, llm)
	return llm, nil
}
