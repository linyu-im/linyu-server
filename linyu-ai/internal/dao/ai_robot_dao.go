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

func (d *aiRobotDao) List(db *gorm.DB) ([]*aiModel.AiRobot, error) {
	var list []*aiModel.AiRobot
	if err := db.Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}
