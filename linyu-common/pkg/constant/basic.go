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
	Friend string // 好友申请
	Group  string // 进群申请
}

// ApplyType 申请类型
var ApplyType = applyType{
	Friend: "friend",
	Group:  "group",
}

func (c applyType) Validate(v string) bool {
	switch v {
	case c.Friend, c.Group:
		return true
	default:
		return false
	}
}

type applySource struct {
	Account string // 账号申请
	ECard   string // 名片申请
	Phone   string // 手机号申请
	Qrcode  string // 二维码申请
}

// ApplySource 申请来源
var ApplySource = applySource{
	Account: "account",
	ECard:   "ecard",
	Phone:   "phone",
	Qrcode:  "qrcode",
}

func (c applySource) Validate(v string) bool {
	switch v {
	case c.Account, c.ECard, c.Phone, c.Qrcode:
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
type contactsPeerType struct {
	Friend     string //好友
	Group      string //群
	Bot        string //机器人
	Enterprise string //企业
}

// ContactsPeerType 通讯录对方的类型
var ContactsPeerType = contactsPeerType{
	Friend:     "friend",
	Group:      "group",
	Bot:        "bot",
	Enterprise: "enterprise",
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
type messageScene struct {
	User  string //用户
	Group string //群
}

// MessageScene 消息场景
var MessageScene = messageScene{
	User:  "user",
	Group: "group",
}

type messageStatus struct {
}

// MessageStatus 消息状态
var MessageStatus = messageStatus{}

type messageType struct {
	Text  string //文本
	Image string //图片
	File  string //文件
	Video string //视频
	Voice string //语音
	ECard string //电子名片
}

// MessageType 消息类型
var MessageType = messageType{
	Text:  "text",
	Image: "image",
	File:  "file",
	Video: "video",
	Voice: "voice",
	ECard: "ecard",
}

type messageFromType struct {
	User string //用户
	Bot  string //机器人
}

// MessageFromType 消息发送放类型
var MessageFromType = messageFromType{
	User: "user",
	Bot:  "bot",
}

// ------------------moment相关------------------

// momentVisibleType 过往可见类型
type momentVisibleType struct {
	All     string //所有人可见
	Include string //指定人可见
	Exclude string //指定人不可见
	Private string //仅自己可见
}

var MomentVisibleType = momentVisibleType{
	All:     "all",
	Include: "include",
	Exclude: "exclude",
	Private: "private",
}

func (m momentVisibleType) Validate(v string) bool {
	switch v {
	case m.All, m.Include, m.Exclude, m.Private:
		return true
	default:
		return false
	}
}

// ------------------enterprise-member相关------------------
type enterpriseMemberRole struct {
	Owner    string // 企业主
	Admin    string // 管理员
	SubAdmin string // 子管理员
	Leader   string // 部门负责人
	Member   string // 普通成员
}

// EnterpriseMemberRole 企业成员身份（Roles 字段可存多个，逗号分隔）
var EnterpriseMemberRole = enterpriseMemberRole{
	Owner:    "owner",
	Admin:    "admin",
	SubAdmin: "subAdmin",
	Leader:   "leader",
	Member:   "member",
}

func (c enterpriseMemberRole) Validate(v string) bool {
	switch v {
	case c.Owner, c.Admin, c.SubAdmin, c.Leader, c.Member:
		return true
	default:
		return false
	}
}
