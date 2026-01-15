package constant

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
