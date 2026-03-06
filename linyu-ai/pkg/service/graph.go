package service

import (
	"context"
	"errors"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
	"sync"
)

var GraphCompilePool sync.Map

func GetRobotGraph(robotId string) (compose.Runnable[map[string]any, *schema.Message], error) {
	if g, ok := GraphCompilePool.Load(robotId); ok {
		return g.(compose.Runnable[map[string]any, *schema.Message]), nil
	}
	//获取机器人信息
	robotInfo := AiRobotService.GetAiRobot(robotId)
	if robotInfo == nil {
		return nil, errors.New("robot not found")
	}
	llm, err := GetLLModel(robotInfo.ModelID)
	if err != nil {
		return nil, err
	}
	graph := compose.NewGraph[map[string]any, *schema.Message]()
	graph.AddLambdaNode("robot_init_node", compose.InvokableLambda(
		func(ctx context.Context, input map[string]any) (map[string]any, error) {
			input["prompt"] = robotInfo.Prompt
			return input, nil
		},
	))
	graph.AddLambdaNode("robot_prompt_node", compose.InvokableLambda(RobotPromptNode))
	graph.AddChatModelNode("llm", llm)

	graph.AddEdge(compose.START, "robot_init_node")
	graph.AddEdge("robot_init_node", "robot_prompt_node")
	graph.AddEdge("robot_prompt_node", "llm")
	graph.AddEdge("llm", compose.END)
	compile, err := graph.Compile(context.Background())
	if err != nil {
		return nil, err
	}
	GraphCompilePool.Store(compile, graph)
	return compile, nil
}
