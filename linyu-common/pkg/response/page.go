package response

import (
	"github.com/linyu-im/linyu-server/linyu-common/pkg/request"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"regexp"
	"strings"
)

type PageResult[T any] struct {
	Records   []T   `json:"records"`
	Total     int64 `json:"total"`
	Page      int   `json:"page"`
	PageSize  int   `json:"pageSize"`
	TotalPage int   `json:"totalPage"`
}

func Paginate[T any](tx *gorm.DB, query request.PageQuery) (*PageResult[T], error) {
	var result PageResult[T]
	var data []T
	var total int64

	query.SetDefault()

	tx.Count(&total)

	if query.SortBy != "" {
		desc := strings.ToUpper(query.SortOrder) == "DESC"
		if tx.Statement != nil && tx.Statement.Schema != nil {
			if _, ok := tx.Statement.Schema.FieldsByDBName[query.SortBy]; ok {
				tx = tx.Order(clause.OrderByColumn{
					Column: clause.Column{Name: query.SortBy},
					Desc:   desc,
				})
			}
		} else {
			re := regexp.MustCompile(`^[a-zA-Z0-9_.]+$`)
			if re.MatchString(query.SortBy) {
				tx = tx.Order(clause.OrderByColumn{
					Column: clause.Column{Name: query.SortBy},
					Desc:   desc,
				})
			}
		}
	}

	offset := (query.Page - 1) * query.PageSize
	if err := tx.Offset(offset).Limit(query.PageSize).Find(&data).Error; err != nil {
		return nil, err
	}

	result.Records = data
	result.Total = total
	result.Page = query.Page
	result.PageSize = query.PageSize
	result.TotalPage = int((total + int64(query.PageSize) - 1) / int64(query.PageSize))

	return &result, nil
}
