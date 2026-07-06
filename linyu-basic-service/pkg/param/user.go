package param

import "github.com/linyu-im/linyu-server/linyu-common/pkg/request"

type UserEmotionSetParam struct {
	EmotionId string `json:"emotionId"`
}

type GetUserAvatarParam struct {
	UserId string `json:"userId"`
}

type UserInfoParam struct {
	UserId string `json:"userId"`
}

type UserUpdateProfileParam struct {
	Username  string `json:"username"`
	Gender    string `json:"gender"`
	Birthday  string `json:"birthday"`
	Signature string `json:"signature"`
	Location  string `json:"location"`
}

type UserSearchParam struct {
	request.PageQuery
	Keyword string `json:"keyword" binding:"required"`
}
