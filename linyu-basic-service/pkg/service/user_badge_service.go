package service

import (
	basicDao "github.com/linyu-im/linyu-server/linyu-basic-service/internal/dao"
	basicModel "github.com/linyu-im/linyu-server/linyu-basic-service/pkg/model"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/constant"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/db"
)

var UserBadgeService = newUserBadgeService()

func newUserBadgeService() *userBadgeService {
	return &userBadgeService{}
}

type userBadgeService struct{}

type badgeConfig struct {
	code      string
	countFunc func(string, string) (int64, error)
}

var badgeConfigs = []badgeConfig{
	{
		code: constant.BadgeCode.NewFriend,
		countFunc: func(uid, lastReadId string) (int64, error) {
			return basicDao.ApplyDao.CountFriendApplyAfterLastReadId(db.RDB, uid, lastReadId)
		},
	},
	{
		code: constant.BadgeCode.GroupNotion,
		countFunc: func(uid, lastReadId string) (int64, error) {
			return basicDao.NoticeDao.CountAfterLastReadId(db.RDB, uid, constant.NoticeType.Group, lastReadId)
		},
	},
}

func (s *userBadgeService) List(userId string) ([]*basicModel.UserBadge, error) {
	list := make([]*basicModel.UserBadge, 0, len(badgeConfigs))
	for _, cfg := range badgeConfigs {
		badge, err := s.getOrCreateBadge(userId, cfg.code)
		if err != nil {
			return nil, err
		}
		count, err := cfg.countFunc(userId, badge.LastReadID)
		if err != nil {
			return nil, err
		}
		badge.UnreadCount = int(count)
		list = append(list, badge)
	}
	return list, nil
}

func (s *userBadgeService) getOrCreateBadge(userId string, badgeCode string) (*basicModel.UserBadge, error) {
	badge := basicDao.UserBadgeDao.GetByUserIdAndCode(db.RDB, userId, badgeCode)
	if badge != nil {
		return badge, nil
	}
	badge = &basicModel.UserBadge{
		UserID:     userId,
		BadgeCode:  badgeCode,
		LastReadID: "0",
	}
	if err := basicDao.UserBadgeDao.Create(db.RDB, badge); err != nil {
		return nil, err
	}
	return badge, nil
}
