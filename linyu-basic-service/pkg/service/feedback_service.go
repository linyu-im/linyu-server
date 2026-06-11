package service

import (
	basicDao "github.com/linyu-im/linyu-server/linyu-basic-service/internal/dao"
	basicModel "github.com/linyu-im/linyu-server/linyu-basic-service/pkg/model"
	basicParam "github.com/linyu-im/linyu-server/linyu-basic-service/pkg/param"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/db"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/utils"
)

var FeedbackService = newFeedbackService()

func newFeedbackService() *feedbackService {
	return &feedbackService{}
}

type feedbackService struct{}

func (s feedbackService) CreateFeedback(userId string, param *basicParam.FeedbackCreateParam) error {
	feedback := &basicModel.Feedback{
		ID:          utils.GenerateSfIDString(),
		UserID:      userId,
		Title:       param.Title,
		Description: param.Description,
		Images:      param.Images,
	}
	if feedback.Images == nil {
		feedback.Images = []string{}
	}
	return basicDao.FeedbackDao.Create(db.RDB, feedback)
}
