package service

import (
	basicDao "github.com/linyu-im/linyu-server/linyu-basic-service/internal/dao"
	basicModel "github.com/linyu-im/linyu-server/linyu-basic-service/pkg/model"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/constant"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/db"
)

var NoticeService = newNoticeService()

func newNoticeService() *noticeService {
	return &noticeService{}
}

type noticeService struct{}

func (s *noticeService) GroupList(userId string) ([]*basicModel.Notice, error) {
	list, err := basicDao.NoticeDao.ListByUserIdAndType(db.RDB, userId, constant.NoticeType.Group)
	if err != nil {
		return nil, err
	}
	lastReadId := "0"
	if len(list) > 0 {
		lastReadId = list[0].ID
	}
	_ = basicDao.UserBadgeDao.UpsertLastReadID(db.RDB, userId, constant.BadgeCode.GroupNotion, lastReadId)
	return list, nil
}
