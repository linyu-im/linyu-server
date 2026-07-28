package constant

// ------------------space-file相关------------------
type spaceType struct {
	User  string
	Group string
	//Org   string
}

var SpaceType = spaceType{
	User:  "user",
	Group: "group",
	//Org:   "org",
}

func (c spaceType) Validate(v string) bool {
	switch v {
	case c.User, c.Group:
		return true
	default:
		return false
	}
}

type spaceStatus struct {
	Active   string
	Disabled string
	Readonly string
}

// SpaceStatus 空间状态
var SpaceStatus = spaceStatus{
	Active:   "active",
	Disabled: "disabled",
	Readonly: "readonly",
}

// DefaultUserSpaceQuotaBytes 用户空间默认容量 50G
const DefaultUserSpaceQuotaBytes int64 = 50 * 1024 * 1024 * 1024

// DefaultSpaceRecycleExpireDays 回收站默认保留天数
const DefaultSpaceRecycleExpireDays = 30

type spaceMemberRole struct {
	Owner  string
	Admin  string
	Member string
}

// SpaceMemberRole 空间成员角色
var SpaceMemberRole = spaceMemberRole{
	Owner:  "owner",
	Admin:  "admin",
	Member: "member",
}

type spaceMemberStatus struct {
	Active   string
	Disabled string
}

// SpaceMemberStatus 空间成员状态
var SpaceMemberStatus = spaceMemberStatus{
	Active:   "active",
	Disabled: "disabled",
}
