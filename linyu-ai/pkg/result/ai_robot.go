package result

import (
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
)

type RobotAnswersPrepareResult struct {
	SessionId string
	Input     map[string]any
	Graph     compose.Runnable[map[string]any, *schema.Message]
}
