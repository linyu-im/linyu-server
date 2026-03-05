package dao

import (
	aiModel "github.com/linyu-im/linyu-server/linyu-ai/pkg/model"
	"gorm.io/gorm"
)

var AiModelDao = newAiModelDao()

func newAiModelDao() *aiModelDao {
	return &aiModelDao{}
}

type aiModelDao struct{}

func (d *aiModelDao) GetById(db *gorm.DB, id string) *aiModel.AiModel {
	result := &aiModel.AiModel{}
	if err := db.First(result, "id = ?", id).Error; err != nil {
		return nil
	}
	return result
}
