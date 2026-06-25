package api

import (
	"github.com/gin-gonic/gin"
	basicParam "github.com/linyu-im/linyu-server/linyu-basic-service/pkg/param"
	basicService "github.com/linyu-im/linyu-server/linyu-basic-service/pkg/service"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/response"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/route"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/utils"
)

func init() {
	route.Register("POST", "/basic/v1/group/create", GroupCreateHandler)
	route.Register("POST", "/basic/v1/group/dissolve", GroupDissolveHandler)
	route.Register("POST", "/basic/v1/group/avatar/get", GetGroupAvatarHandler)
	route.Register("POST", "/basic/v1/group/invite-member", GroupInviteMemberHandler)
	route.Register("POST", "/basic/v1/group/remove-member", GroupRemoveMemberHandler)
	route.Register("POST", "/basic/v1/group/info", GroupInfoHandler)
	route.Register("POST", "/basic/v1/group/member/list", GroupMemberListHandler)
}

// GroupCreateHandler 群聊创建
func GroupCreateHandler(c *gin.Context) {
	param := &basicParam.GroupCreateParam{}
	if !utils.ShouldBindBodyWithJSONAndValidate(c, param) {
		return
	}
	currentUserId := c.GetString("userId")
	err := basicService.GroupService.GroupCreate(currentUserId, param)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Ok(c)
}

// GroupDissolveHandler 群聊解散
func GroupDissolveHandler(c *gin.Context) {
	param := &basicParam.GroupDissolveParam{}
	if !utils.ShouldBindBodyWithJSONAndValidate(c, param) {
		return
	}
	currentUserId := c.GetString("userId")
	err := basicService.GroupService.GroupDissolve(currentUserId, param)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Ok(c)
}

// GroupInviteMemberHandler 群聊邀请成员
func GroupInviteMemberHandler(c *gin.Context) {
	param := &basicParam.GroupInviteMemberParam{}
	if !utils.ShouldBindBodyWithJSONAndValidate(c, param) {
		return
	}
	currentUserId := c.GetString("userId")
	err := basicService.GroupService.InviteMember(currentUserId, param)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Ok(c)
}

// GroupRemoveMemberHandler 群聊移除成员
func GroupRemoveMemberHandler(c *gin.Context) {
	param := &basicParam.GroupRemoveMemberParam{}
	if !utils.ShouldBindBodyWithJSONAndValidate(c, param) {
		return
	}
	currentUserId := c.GetString("userId")
	err := basicService.GroupService.RemoveMember(currentUserId, param)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Ok(c)
}

func GetGroupAvatarHandler(c *gin.Context) {
	param := &basicParam.GetGroupAvatarParam{}
	if !utils.ShouldBindBodyWithJSONAndValidate(c, param) {
		return
	}
	url := basicService.GroupService.GetGroupAvatar(param.GroupId)
	response.Ok(c, url)
}

func GroupInfoHandler(c *gin.Context) {
	param := &basicParam.GroupInfoParam{}
	if !utils.ShouldBindBodyWithJSONAndValidate(c, param) {
		return
	}
	currentUserId := c.GetString("userId")
	group := basicService.GroupService.GroupInfo(currentUserId, param.GroupId)
	response.Ok(c, group)
}

// GroupMemberListHandler 获取群成员列表
func GroupMemberListHandler(c *gin.Context) {
	param := &basicParam.GroupMemberListParam{}
	if !utils.ShouldBindBodyWithJSONAndValidate(c, param) {
		return
	}
	currentUserId := c.GetString("userId")
	data, err := basicService.GroupService.GroupMemberList(currentUserId, param)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Ok(c, data)
}
