package constant

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
