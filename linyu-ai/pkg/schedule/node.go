package schedule

import (
	"context"
	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/schema"
)

func RobotPromptNode(ctx context.Context, input map[string]any) ([]*schema.Message, error) {
	//短期记忆
	messages, ok := input["shortMemory"].([]*schema.Message)
	if !ok {
		messages = []*schema.Message{}
	}
	template := prompt.FromMessages(schema.FString,
		// 系统消息模板
		schema.SystemMessage(input["prompt"].(string)),
		// 用户消息模板
		schema.UserMessage("问题:{question}"),
	)
	msg, _ := template.Format(context.Background(), map[string]any{
		"question": input["question"],
	})
	messages = append(messages, msg...)
	return messages, nil
}
