package model

import (
	"encoding/json"
	"fmt"

	"github.com/linyu-im/linyu-server/linyu-common/pkg/constant"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/db"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/localtime"
	"gorm.io/gorm"
)

func init() {
	db.AddMigrateTable(&Message{})
}

// Message 消息表
type Message struct {
	ID             string              `gorm:"size:64;primaryKey;autoIncrement:false;comment:消息id" json:"id"`
	SessionID      string              `gorm:"size:256;index;not null;comment:会话id" json:"sessionId"`
	FromID         string              `gorm:"size:64;index;not null;comment:消息发送方id" json:"fromId"`
	ToID           string              `gorm:"size:64;index;not null;comment:消息接受方id" json:"toId"`
	MsgType        string              `gorm:"size:64;comment:消息类型" json:"msgType"`
	FromType       string              `gorm:"size:64;comment:发送方类型" json:"fromType"`
	IsShowTime     bool                `gorm:"comment:是否显示时间;default:0" json:"isShowTime"`
	Content        json.RawMessage     `gorm:"type:text;serializer:json;comment:消息内容" json:"content"`
	KeywordContent string              `gorm:"type:text;comment:关键内容" json:"keywordContent"`
	Status         string              `gorm:"size:64;comment:消息状态" json:"status"`
	SceneType      string              `gorm:"size:64;not null;comment:会话场景" json:"sceneType"`
	QuoteMsgId     string              `gorm:"size:64;comment:引用消息的id" json:"quoteMsgId"`
	CreatedAt      localtime.LocalTime `gorm:"type:timestamp(3);not null;autoCreateTime;comment:创建时间" json:"createdAt"`
	UpdatedAt      localtime.LocalTime `gorm:"type:timestamp(3);not null;autoUpdateTime;comment:更新时间" json:"updatedAt"`
	DeletedAt      gorm.DeletedAt      `gorm:"index" json:"deletedAt"`
}

func (Message) TableName() string {
	return "t_message"
}

func (Message) TableComment() string {
	return "消息表"
}

var msgContentFactory = map[string]func() MsgContent{
	constant.MessageType.Text:    func() MsgContent { return &TextContent{} },
	constant.MessageType.Image:   func() MsgContent { return &ImageContent{} },
	constant.MessageType.Video:   func() MsgContent { return &VideoContent{} },
	constant.MessageType.File:    func() MsgContent { return &FileContent{} },
	constant.MessageType.ECard:   func() MsgContent { return &ECardContent{} },
	constant.MessageType.Voice:   func() MsgContent { return &VoiceContent{} },
	constant.MessageType.Sticker: func() MsgContent { return &StickerContent{} },
}

func ParseMsgContent(msgType string, raw json.RawMessage) (MsgContent, error) {
	factory, ok := msgContentFactory[msgType]
	if !ok {
		return nil, fmt.Errorf("unsupported msgType: %s", msgType)
	}
	content := factory()
	if err := json.Unmarshal(raw, content); err != nil {
		return nil, err
	}
	return content, nil
}

func MsgContentToString(msgType string, raw json.RawMessage) (string, error) {
	content, err := ParseMsgContent(msgType, raw)
	if err != nil {
		return "", err
	}
	return content.ToString(), nil
}

type MsgContent interface {
	ToString() string
	GetKeywordContent() string
}

type TextContent struct {
	Text string `json:"text"`
}

func (c TextContent) ToString() string {
	return c.Text
}

func (c TextContent) GetKeywordContent() string {
	return c.Text
}

type ImageContent struct {
	ImgUrl      string `json:"imgUrl"`
	ImgThumbUrl string `json:"imgThumbUrl"`
	ImgName     string `json:"imgName"`
	ImgSize     int64  `json:"imgSize"`
}

func (c ImageContent) ToString() string {
	return "[Image URL]:" + c.ImgUrl
}

func (c ImageContent) GetKeywordContent() string {
	return ""
}

type VideoContent struct {
	VideoUrl      string `json:"videoUrl"`
	VideoThumbUrl string `json:"VideoThumbUrl"`
	VideoName     string `json:"videoName"`
	VideoSize     int64  `json:"videoSize"`
}

func (c VideoContent) ToString() string {
	return "[Video URL]:" + c.VideoUrl
}

func (c VideoContent) GetKeywordContent() string {
	return ""
}

type FileContent struct {
	FileUrl  string `json:"fileUrl"`
	FileType string `json:"fileType"`
	FileName string `json:"fileName"`
	FileSize int64  `json:"fileSize"`
}

func (c FileContent) ToString() string {
	return "[File URL]:" + c.FileUrl
}

func (c FileContent) GetKeywordContent() string {
	return c.FileName
}

type ECardContent struct {
	UserID   string `json:"userId"`
	UserName string `json:"userName"`
}

func (c ECardContent) ToString() string {
	return "[E-Card]:" + c.UserName
}

func (c ECardContent) GetKeywordContent() string {
	return c.UserName
}

type VoiceContent struct {
	VoiceUrl      string `json:"voiceUrl"`
	VoiceDuration string `json:"voiceDuration"`
}

func (c VoiceContent) ToString() string {
	return "[Voice URL]:" + c.VoiceUrl
}

func (c VoiceContent) GetKeywordContent() string {
	return ""
}

type StickerContent struct {
	StickerUrl  string `json:"stickerUrl"`
	StickerName string `json:"stickerName"`
}

func (c StickerContent) ToString() string {
	return "[Sticker URL]:" + c.StickerUrl
}

func (c StickerContent) GetKeywordContent() string {
	return ""
}
