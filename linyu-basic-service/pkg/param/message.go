package param

import (
	"encoding/json"
)

type SendMessageToUserParam struct {
	ToUserId string          `json:"toUserId" binding:"required"`
	MsgType  string          `json:"msgType" binding:"required"`
	Content  json.RawMessage `json:"content" binding:"required"`
}

type SendMessageToGroupParam struct {
	ToGroupId string          `json:"toGroupId" binding:"required"`
	MsgType   string          `json:"msgType" binding:"required"`
	Content   json.RawMessage `json:"content" binding:"required"`
}
