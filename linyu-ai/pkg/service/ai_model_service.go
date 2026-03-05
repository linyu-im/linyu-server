package service

import (
	aiDao "github.com/linyu-im/linyu-server/linyu-ai/internal/dao"
	aiModel "github.com/linyu-im/linyu-server/linyu-ai/pkg/model"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/db"
)

var AiModelService = newAiModelService()

func newAiModelService() *aiModelService {
	return &aiModelService{}
}

type aiModelService struct{}

func (s *aiModelService) GetAiModel(modelId string) *aiModel.AiModel {
	model := aiDao.AiModelDao.GetById(db.RDB, modelId)
	return model
}
