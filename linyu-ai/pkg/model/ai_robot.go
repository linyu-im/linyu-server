package model

import (
	"github.com/linyu-im/linyu-server/linyu-common/pkg/db"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/localtime"
	"gorm.io/gorm"
)

func init() {
	db.AddMigrateTable(&AiRobot{})
}

// AiRobot AI机器人表
type AiRobot struct {
	ID          string              `gorm:"size:64;primaryKey;autoIncrement:false;comment:id" json:"id"`
	RobotName   string              `gorm:"size:255;not null;comment:机器人名称" json:"robotName"`
	RobotAvatar string              `gorm:"size:512;comment:机器人头像" json:"robotAvatar"`
	RobotDesc   string              `gorm:"size:1024;comment:机器人描述" json:"robotDesc"`
	ModelID     string              `gorm:"size:64;comment:模型id" json:"modelId"`
	Prompt      string              `gorm:"type:longtext;comment:机器人提示词" json:"prompt"`
	Status      string              `gorm:"size:64;default:'active';comment:机器人状态" json:"status"`
	CreatedAt   localtime.LocalTime `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt   localtime.LocalTime `gorm:"autoUpdateTime" json:"updatedAt"`
	DeletedAt   gorm.DeletedAt      `gorm:"index" json:"deletedAt"`
}

func (AiRobot) TableName() string {
	return "t_ai_robot"
}

func (AiRobot) TableComment() string {
	return "AI机器人表"
}
