package service

import (
	aiDao "github.com/linyu-im/linyu-server/linyu-ai/internal/dao"
	aiModel "github.com/linyu-im/linyu-server/linyu-ai/pkg/model"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/db"
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
