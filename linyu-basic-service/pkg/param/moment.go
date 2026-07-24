package param

import (
	"github.com/linyu-im/linyu-server/linyu-basic-service/pkg/model"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/request"
)

type MomentCreateParam struct {
	TextContent    string                `json:"textContent"`
	MediaContent   []*model.MediaContent `json:"mediaContent"`
	VisibleType    string                `json:"visibleType"`
	VisibleUserIds []string              `json:"visibleUserIds"`
	Location       string                `json:"location"`
}

type MomentPageParam struct {
	request.PageQuery
	ViewUserId string `json:"viewUserId"`
}

type MomentLikeAddParam struct {
	MomentId string `json:"momentId" binding:"required"`
}

type MomentLikeCancelParam struct {
	MomentId string `json:"momentId" binding:"required"`
}

type MomentCommentAddParam struct {
	MomentId string `json:"momentId" binding:"required"`
	ParentID string `json:"parentId"`
	Content  string `json:"content" binding:"required"`
}

type MomentCommentDelParam struct {
	CommentId string `json:"commentId" binding:"required"`
}

type MomentDeleteParam struct {
	MomentId string `json:"momentId" binding:"required"`
}

type MomentSettingParam struct {
	UserId string `json:"userId" binding:"required"`
}
