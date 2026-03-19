package model

import (
	"github.com/linyu-im/linyu-server/linyu-common/pkg/db"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/localtime"
	"gorm.io/gorm"
)

func init() {
	db.MysqlAddMigrateTable(&Message{})
}

// Message 消息表
type Message struct {
	ID         string              `gorm:"size:64;primaryKey;autoIncrement:false;comment:消息id" json:"id"`
	SessionID  string              `gorm:"size:256;index;not null;comment:会话id" json:"sessionId"`
	FromID     string              `gorm:"size:64;index;not null;comment:消息发送方id" json:"fromId"`
	ToID       string              `gorm:"size:64;index;not null;comment:消息接受方id" json:"toId"`
	MsgType    string              `gorm:"size:64;comment:消息类型" json:"msgType"`
	FromType   string              `gorm:"size:64;comment:发送方类型" json:"fromType"`
	IsShowTime bool                `gorm:"comment:是否显示时间;default:0" json:"isShowTime"`
	Content    MsgContent          `gorm:"type:text;serializer:json;comment:消息内容" json:"content"`
	Status     string              `gorm:"size:64;comment:消息状态" json:"status"`
	MsgScene   string              `gorm:"size:64;not null;comment:消息场景" json:"MsgScene"`
	QuoteMsgId string              `gorm:"size:64;comment:引用消息的id" json:"quoteMsgId"`
	CreatedAt  localtime.LocalTime `gorm:"type:timestamp(3);not null;autoCreateTime;comment:创建时间" json:"createdAt"`
	UpdatedAt  localtime.LocalTime `gorm:"type:timestamp(3);not null;autoUpdateTime;comment:更新时间" json:"updatedAt"`
	DeletedAt  gorm.DeletedAt      `gorm:"index" json:"deletedAt"`
}

func (Message) TableName() string {
	return "t_message"
}

func (Message) TableComment() string {
	return "消息表"
}

type MsgContent interface {
	ToString() string
}

type TextContent struct {
	Text string `json:"text"`
}

func (c TextContent) ToString() string {
	return c.Text
}

type ImageContent struct {
	ImgUrl      string `json:"imgUrl"`
	ImgThumbUrl string `json:"imgThumbUrl"`
	ImgName     string `json:"imgName"`
	ImgSize     string `json:"imgSize"`
}

func (c ImageContent) ToString() string {
	return "[Image URL]:" + c.ImgUrl
}

type VideoContent struct {
	VideoUrl      string `json:"videoUrl"`
	VideoThumbUrl string `json:"VideoThumbUrl"`
	VideoName     string `json:"videoName"`
	VideoSize     string `json:"videoSize"`
}

func (c VideoContent) ToString() string {
	return "[Video URL]:" + c.VideoUrl
}

type FileContent struct {
	FileUrl  string `json:"fileUrl"`
	FileType string `json:"fileType"`
	FileName string `json:"fileName"`
	FileSize string `json:"fileSize"`
}

func (c FileContent) ToString() string {
	return "[File URL]:" + c.FileUrl
}

type ECardContent struct {
	UserID     string `json:"userId"`
	UserName   string `json:"userName"`
	UserAvatar string `json:"userAvatar"`
}

func (c ECardContent) ToString() string {
	return "[E-Card]:" + c.UserName + ":" + c.UserAvatar
}

type VoiceContent struct {
	VoiceUrl      string `json:"voiceUrl"`
	VoiceDuration string `json:"voiceDuration"`
}

func (c VoiceContent) ToString() string {
	return "[Voice URL]:" + c.VoiceUrl
}
