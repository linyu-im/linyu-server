package param

type GroupCreateParam struct {
	GroupName       string   `json:"groupName" binding:"required"`
	GroupMemberList []string `json:"groupMemberList"`
}

type GroupDissolveParam struct {
	GroupId string `json:"groupId" binding:"required"`
}
