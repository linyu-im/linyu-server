package constant

type messageSource struct {
	User  string //用户
	Group string //群
}

// MessageSource 消息源
var MessageSource = messageSource{
	User:  "user",
	Group: "group",
}

type messageStatus struct {
	Read   string //已读
	Unread string //未读
}

// MessageStatus 消息状态
var MessageStatus = messageStatus{
	Read:   "read",
	Unread: "unread",
}

type messageType struct {
	Text  string //文本
	Image string //图片
	File  string //文件
}

// MessageType 消息类型
var MessageType = messageType{
	Text:  "text",
	Image: "image",
	File:  "file",
}
