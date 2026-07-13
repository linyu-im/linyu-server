package param

import "github.com/linyu-im/linyu-server/linyu-common/pkg/request"

type GroupCreateParam struct {
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

type GroupUpdateInfoParam struct {
	GroupId  string `json:"groupId" binding:"required"`
	Name     string `json:"name"`
	Describe string `json:"describe"`
	Tag      string `json:"tag"`
}

type GroupSetAdminParam struct {
	GroupId         string   `json:"groupId" binding:"required"`
	AddAdminList    []string `json:"addAdminList"`
	RemoveAdminList []string `json:"removeAdminList"`
}

type GroupTransferOwnerParam struct {
	GroupId    string `json:"groupId" binding:"required"`
	NewOwnerId string `json:"newOwnerId" binding:"required"`
}

type GroupLeaveParam struct {
	GroupId string `json:"groupId" binding:"required"`
}

type GroupNoticeListParam struct {
	GroupId string `json:"groupId" binding:"required"`
}

type GroupNoticeAddParam struct {
	GroupId string `json:"groupId" binding:"required"`
	Content string `json:"content" binding:"required"`
	IsTop   bool   `json:"isTop"`
}

type GroupNoticeUpdateParam struct {
	NoticeId string `json:"noticeId" binding:"required"`
	Content  string `json:"content"`
	IsTop    *bool  `json:"isTop"`
}

type GroupNoticeDeleteParam struct {
	NoticeId string `json:"noticeId" binding:"required"`
}
