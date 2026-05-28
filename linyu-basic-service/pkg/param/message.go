package param

import (
	"encoding/json"
)

type SendMessageToUserParam struct {
	ToUserId string          `json:"toUserId" binding:"required"`
	MsgType  string          `json:"msgType" binding:"required"`
	Content  json.RawMessage `json:"content" binding:"required"`
	Mentions []struct {
		Id          string `json:"id"`
		MentionType string `json:"mentionType"`
	} `json:"mentions"`
}

type SendMessageToGroupParam struct {
	ToGroupId string          `json:"toGroupId" binding:"required"`
	MsgType   string          `json:"msgType" binding:"required"`
	Content   json.RawMessage `json:"content" binding:"required"`
	Mentions  []struct {
		Id          string `json:"id"`
		MentionType string `json:"type"`
	} `json:"mentions"`
}

type MessagePageParam struct {
	ToId     string `json:"toId" binding:"required"`
	Page     int    `json:"page"`
	PageSize int    `json:"pageSize"`
}

type UploadMsgFileInfoParam struct {
	FileHash   string `json:"fileHash" binding:"required"`
	FileSize   int64  `json:"fileSize" binding:"required"`
	FileName   string `json:"fileName" binding:"required"`
	TotalChunk int    `json:"totalChunk" binding:"required"`
}
