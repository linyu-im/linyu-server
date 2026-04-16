package dao

import (
	basicModel "github.com/linyu-im/linyu-server/linyu-basic-service/pkg/model"
	"gorm.io/gorm"
)

var EmotionDao = newEmotionDao()

func newEmotionDao() *emotionDao {
	return &emotionDao{}
}

type emotionDao struct{}

func (d *emotionDao) EmotionList(db *gorm.DB) ([]*basicModel.Emotion, error) {
	var list []*basicModel.Emotion
	if err := db.Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (d *emotionDao) SetEmotion(db *gorm.DB, userId string, emotionId string) error {
	return db.
		Model(&basicModel.User{}).
		Where("id = ?", userId).
		Update("emotion_id", emotionId).Error
}
