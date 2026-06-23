package service

import (
	aiDao "github.com/linyu-im/linyu-server/linyu-ai/internal/dao"
	aiModel "github.com/linyu-im/linyu-server/linyu-ai/pkg/model"
	aiParam "github.com/linyu-im/linyu-server/linyu-ai/pkg/param"
	aiResult "github.com/linyu-im/linyu-server/linyu-ai/pkg/result"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/constant"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/db"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/utils"
)

var AiRobotService = newAiRobotService()

func newAiRobotService() *aiRobotService {
	return &aiRobotService{}
}

type aiRobotService struct{}

func (s *aiRobotService) GetAiRobot(robotId string) *aiModel.AiRobot {
	robot := aiDao.AiRobotDao.GetById(db.RDB, robotId)
	return robot
}

func (s *aiRobotService) List() ([]*aiModel.AiRobot, error) {
	return aiDao.AiRobotDao.List(db.RDB)
}

func (s *aiRobotService) GetAvatar(robotId string) string {
	robot := aiDao.AiRobotDao.GetById(db.RDB, robotId)
	if robot == nil {
		return ""
	}
	return robot.RobotAvatar
}

func (s *aiRobotService) PrepareRobotAnswers(userId string, param *aiParam.RobotAnswersParam) (*aiResult.RobotAnswersPrepareResult, error) {
	var sessionId string
	switch param.SceneType {
	case constant.SceneType.User:
		sessionId = utils.Generate1v1SessionID(userId, param.PeerId)
	case constant.SceneType.Group:
		sessionId = param.PeerId
	}
	//longMemory, _ := AiMemoryService.GetLongTermMemory(sessionId, param.Question)
	shortMemory := AiMemoryService.GetShortTermMemory(sessionId)
	in := map[string]any{
		"question":    param.Question,
		"robotId":     param.RobotId,
		"sessionId":   sessionId,
		"shortMemory": shortMemory,
		//"longMemory":  longMemory,
	}
	graph, err := GetRobotGraph(param.RobotId)
	if err != nil {
		return nil, err
	}
	return &aiResult.RobotAnswersPrepareResult{
		SessionId: sessionId,
		Input:     in,
		Graph:     graph,
	}, nil
}
