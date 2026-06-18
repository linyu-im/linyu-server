package service

import (
	"context"
	"fmt"
	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/schema"
)

func RobotPromptNode(ctx context.Context, input map[string]any) ([]*schema.Message, error) {
	var messages []*schema.Message
	// 系统提示词
	template := prompt.FromMessages(schema.FString,
		schema.SystemMessage(input["prompt"].(string)),
	)
	msg, _ := template.Format(context.Background(), map[string]any{
		"question": input["question"],
	})
	messages = append(messages, msg...)
	//长期记忆
	if longMemory, ok := input["longMemory"].(*schema.Message); ok && longMemory != nil {
		messages = append(messages, longMemory)
	}
	//短期记忆
	if shortMemory, ok := input["shortMemory"].([]*schema.Message); ok && len(shortMemory) > 0 {
		messages = append(messages, shortMemory...)
	}
	// 用户消息
	messages = append(messages, schema.UserMessage(fmt.Sprintf("问题:%v", input["question"])))
	return messages, nil
}
