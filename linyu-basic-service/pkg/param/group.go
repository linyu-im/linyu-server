package param

import "github.com/linyu-im/linyu-server/linyu-common/pkg/request"

type GroupCreateParam struct {
	GroupName       string   `json:"groupName" binding:"required"`
	GroupMemberList []string `json:"groupMemberList"`
}

type GroupDissolveParam struct {
	GroupId string `json:"groupId" binding:"required"`
}

type GroupInviteMemberParam struct {
	GroupId         string   `json:"groupId" binding:"required"`
	GroupMemberList []string `json:"groupMemberList"`
}

type GroupRemoveMemberParam struct {
	GroupId         string   `json:"groupId" binding:"required"`
	GroupMemberList []string `json:"groupMemberList"`
}

type GetGroupAvatarParam struct {
	GroupId string `json:"groupId" binding:"required"`
}

type GroupInfoParam struct {
	GroupId string `json:"groupId" binding:"required"`
}

type GroupMemberListParam struct {
	GroupId string `json:"groupId" binding:"required"`
}

type GroupSearchParam struct {
	request.PageQuery
	Keyword string `json:"keyword" binding:"required"`
}
