package dao

import (
	aiModel "github.com/linyu-im/linyu-server/linyu-ai/pkg/model"
	"gorm.io/gorm"
)

var AiRobotDao = newAiRobotDao()

func newAiRobotDao() *aiRobotDao {
	return &aiRobotDao{}
}

type aiRobotDao struct{}

func (d *aiRobotDao) GetById(db *gorm.DB, id string) *aiModel.AiRobot {
	result := &aiModel.AiRobot{}
	if err := db.First(result, "id = ?", id).Error; err != nil {
		return nil
	}
	return result
}
