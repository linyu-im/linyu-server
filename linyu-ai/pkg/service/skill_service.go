package service

import (
	aiDao "github.com/linyu-im/linyu-server/linyu-ai/internal/dao"
	aiModel "github.com/linyu-im/linyu-server/linyu-ai/pkg/model"
	aiParam "github.com/linyu-im/linyu-server/linyu-ai/pkg/param"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/db"
)

var AiSkillService = newAiSkillService()

func newAiSkillService() *aiSkillService {
	return &aiSkillService{}
}

type aiSkillService struct{}

func (s *aiSkillService) List(param *aiParam.SkillListParam) ([]*aiModel.Skill, error) {
	return aiDao.SkillDao.List(db.RDB, param)
}
