package service

import (
	"errors"
	basicDao "github.com/linyu-im/linyu-server/linyu-basic-service/internal/dao"
	basicModel "github.com/linyu-im/linyu-server/linyu-basic-service/pkg/model"
	"github.com/linyu-im/linyu-server/linyu-basic-service/pkg/param"
	basicResult "github.com/linyu-im/linyu-server/linyu-basic-service/pkg/result"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/constant"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/db"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/response"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/utils"
	"gorm.io/gorm"
)

var MomentService = newMomentService()

func newMomentService() *momentService {
	return &momentService{}
}

type momentService struct{}

func (s momentService) CreateMoment(userId string, param *param.MomentCreateParam) error {
	if !constant.MomentVisibleType.Validate(param.VisibleType) {
		return errors.New("param.type-error")
	}
	moment := &basicModel.Moment{
		ID:            utils.GenerateSfIDString(),
		UserID:        userId,
		VisibleType:   param.VisibleType,
		MediaContents: param.MediaContent,
		TextContent:   param.TextContent,
		Location:      param.Location,
	}
	err := db.RDB.Transaction(func(tx *gorm.DB) error {
		//创建过往可见关系
		if constant.MomentVisibleType.Include == param.VisibleType ||
			constant.MomentVisibleType.Exclude == param.VisibleType {
			var userVisibles []*basicModel.MomentVisible
			for _, id := range param.VisibleUserIds {
				userVisibles = append(userVisibles, &basicModel.MomentVisible{
					ID:          utils.GenerateSfIDString(),
					MomentID:    moment.ID,
					UserID:      id,
					VisibleType: param.VisibleType,
				})
			}
			if err := basicDao.MomentVisibleDao.CreateBatch(tx, userVisibles); err != nil {
				return err
			}
		}
		//创建过往
		if err := basicDao.MomentDao.Create(tx, moment); err != nil {
			return err
		}
		return nil
	})
	return err
}

func (s momentService) PageMoment(
	userId string,
	param *param.MomentPageParam,
) (*response.PageResult[*basicResult.MomentPageResult], error) {
	// 构建过往查询相关条件
	tx := basicDao.MomentDao.BuildMomentQuery(db.RDB, userId, param.ViewUserId)
	pages, err := response.Paginate[*basicModel.Moment](tx, param.PageQuery)
	if err != nil {
		return nil, err
	}

	var records = make([]*basicResult.MomentPageResult, 0)
	for _, moment := range pages.Records {
		comments := basicDao.MomentCommentDao.ListByMomentId(db.RDB, moment.ID)
		likes := basicDao.MomentLikeDao.ListByMomentId(db.RDB, moment.ID)
		records = append(records, &basicResult.MomentPageResult{
			Moment:   moment,
			Comments: comments,
			Likes:    likes,
		})
	}

	result := &response.PageResult[*basicResult.MomentPageResult]{
		Records:   records,
		Page:      pages.Page,
		PageSize:  param.PageSize,
		Total:     pages.Total,
		TotalPage: pages.TotalPage,
	}
	return result, nil
}
