package result

import basicModel "github.com/linyu-im/linyu-server/linyu-basic-service/pkg/model"

type MomentPageResult struct {
	Moment   *basicModel.Moment          `json:"moment"`
	Comments []*basicModel.MomentComment `json:"comments"`
	Likes    []*basicModel.MomentLike    `json:"likes"`
}
