package constant

type memberRole struct {
	General       string
	Administrator string
}

// MemberRole 群成员角色
var MemberRole = memberRole{
	General:       "general",
	Administrator: "administrator",
}
