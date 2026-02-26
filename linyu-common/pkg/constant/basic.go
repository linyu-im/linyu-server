package constant

// ------------------space-file相关------------------
type userStatus struct {
	Active string
	Banned string
}

var UserStatus = userStatus{
	Active: "active",
	Banned: "banned",
}

// ------------------apply相关------------------
type applyType struct {
	AddContacts string // 添加联系人申请
	JoinGroup   string // 进群申请
}

// ApplyType 申请类型
var ApplyType = applyType{
	AddContacts: "addContacts",
	JoinGroup:   "joinGroup",
}

func (c applyType) Validate(v string) bool {
	switch v {
	case c.AddContacts, c.JoinGroup:
		return true
	default:
		return false
	}
}

type applyStatus struct {
	Wait   string //等待
	Agree  string //同意
	Cancel string //取消
	Reject string //拒绝
}

// ApplyStatus 申请状态
var ApplyStatus = applyStatus{
	Wait:   "wait",
	Agree:  "agree",
	Cancel: "cancel",
	Reject: "reject",
}

// ------------------chat相关------------------
type chatType struct {
	User  string // 用户
	Group string // 群
	Bot   string // 机器人
}

// ChatType 聊天会话类型
var ChatType = chatType{
	User:  "user",
	Group: "group",
	Bot:   "bot",
}

func (c chatType) Validate(v string) bool {
	switch v {
	case c.User, c.Group, c.Bot:
		return true
	default:
		return false
	}
}

// ------------------contacts相关------------------
type contactsType struct {
	User  string //用户
	Group string //群
	Bot   string //机器人
}

// ContactsType 通讯录数据类型
var ContactsType = contactsType{
	User:  "user",
	Group: "group",
	Bot:   "bot",
}

// ------------------group-member相关------------------
type memberRole struct {
	General       string
	Administrator string
}

// MemberRole 群成员角色
var MemberRole = memberRole{
	General:       "general",
	Administrator: "administrator",
}

// ------------------message相关------------------
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
