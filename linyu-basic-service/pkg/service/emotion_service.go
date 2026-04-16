package service

import (
	basicDao "github.com/linyu-im/linyu-server/linyu-basic-service/internal/dao"
	basicModel "github.com/linyu-im/linyu-server/linyu-basic-service/pkg/model"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/db"
)

var EmotionService = newEmotionService()

func newEmotionService() *emotionService {
	return &emotionService{}
}

type emotionService struct{}

func (s emotionService) EmotionList() ([]*basicModel.Emotion, error) {
	return basicDao.EmotionDao.EmotionList(db.RDB)
}

func (s emotionService) SetEmotion(userId string, emotionId string) error {
	return basicDao.EmotionDao.SetEmotion(db.RDB, userId, emotionId)
}
