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
