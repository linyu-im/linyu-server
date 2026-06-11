package dao

import (
	basicModel "github.com/linyu-im/linyu-server/linyu-basic-service/pkg/model"
	"gorm.io/gorm"
)

var FeedbackDao = newFeedbackDao()

func newFeedbackDao() *feedbackDao {
	return &feedbackDao{}
}

type feedbackDao struct{}

func (d *feedbackDao) Create(db *gorm.DB, feedback *basicModel.Feedback) error {
	return db.Create(feedback).Error
}
