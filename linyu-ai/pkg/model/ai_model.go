package model

import (
	"github.com/linyu-im/linyu-server/linyu-common/pkg/db"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/localtime"
	"gorm.io/gorm"
)

func init() {
	db.AddMigrateTable(&AiModel{})
}

// AiModel AI模型表
type AiModel struct {
	ID              string              `gorm:"size:64;primaryKey;autoIncrement:false;comment:id" json:"id"`
	ModelName       string              `gorm:"size:255;not null;comment:模型名称" json:"name"`
	ModelLogo       string              `gorm:"size:512;comment:模型logo" json:"logo"`
	ModelDesc       string              `gorm:"size:1024;comment:模型的描述" json:"desc"`
	BaseURL         string              `gorm:"size:512;comment:模型的地址" json:"baseUrl"`
	Model           string              `gorm:"size:255;comment:模型" json:"model"`
	APIKey          string              `gorm:"size:255;comment:模型ApiKey" json:"apiKey"`
	MaxTokens       *int                `gorm:"comment:模型最大token数" json:"maxTokens"`
	Temperature     *float32            `gorm:"comment:模型温度" json:"temperature"`
	ReasoningEffort string              `gorm:"size:255;comment:推理努力" json:"reasoningEffort"`
	Status          string              `gorm:"size:64;default:'active';comment:机器人状态" json:"status"`
	CreatedAt       localtime.LocalTime `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt       localtime.LocalTime `gorm:"autoUpdateTime" json:"updatedAt"`
	DeletedAt       gorm.DeletedAt      `gorm:"index" json:"deletedAt"`
}

func (AiModel) TableName() string {
	return "t_ai_model"
}

func (AiModel) TableComment() string {
	return "AI模型表"
}
