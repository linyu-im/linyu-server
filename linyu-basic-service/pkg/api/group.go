package api

import (
	"fmt"
	"path/filepath"

	"github.com/gin-gonic/gin"
	basicParam "github.com/linyu-im/linyu-server/linyu-basic-service/pkg/param"
	basicService "github.com/linyu-im/linyu-server/linyu-basic-service/pkg/service"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/response"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/route"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/storage"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/utils"
)

func init() {
	route.Register("POST", "/basic/v1/group/create", GroupCreateHandler)
	route.Register("POST", "/basic/v1/group/dissolve", GroupDissolveHandler)
	route.Register("POST", "/basic/v1/group/avatar/get", GetGroupAvatarHandler)
	route.Register("POST", "/basic/v1/group/avatar/upload", UploadGroupAvatarHandler)
	route.Register("POST", "/basic/v1/group/invite-member", GroupInviteMemberHandler)
	route.Register("POST", "/basic/v1/group/remove-member", GroupRemoveMemberHandler)
	route.Register("POST", "/basic/v1/group/info", GroupInfoHandler)
	route.Register("POST", "/basic/v1/group/member/list", GroupMemberListHandler)
	route.Register("POST", "/basic/v1/group/search", GroupSearchHandler)
}

// GroupCreateHandler 群聊创建
func GroupCreateHandler(c *gin.Context) {
	param := &basicParam.GroupCreateParam{}
	if !utils.ShouldBindBodyWithJSONAndValidate(c, param) {
		return
	}
	currentUserId := c.GetString("userId")
	group, err := basicService.GroupService.GroupCreate(currentUserId, param)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Ok(c, group)
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

// UploadGroupAvatarHandler 上传群头像
func UploadGroupAvatarHandler(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		response.Fail(c, "param.file-not-found")
		return
	}

	groupId := c.PostForm("groupId")
	if groupId == "" {
		response.Fail(c, "param.error")
		return
	}

	const maxFileSize = 10 * 1024 * 1024
	if file.Size > maxFileSize {
		response.Fail(c, "file too large, max 10MB")
		return
	}

	currentUserId := c.GetString("userId")
	ext := filepath.Ext(file.Filename)
	fileKey := fmt.Sprintf("avatar/group/%s%s", groupId, ext)

	src, err := file.Open()
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	defer src.Close()

	url, err := storage.S.Upload(fileKey, src)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}

	if err := basicService.GroupService.UpdateAvatar(currentUserId, groupId, url); err != nil {
		response.Fail(c, err.Error())
		return
	}
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

// GroupSearchHandler 模糊查询群聊列表
func GroupSearchHandler(c *gin.Context) {
	param := &basicParam.GroupSearchParam{}
	if !utils.ShouldBindBodyWithJSONAndValidate(c, param) {
		return
	}
	result, err := basicService.GroupService.SearchByKeyword(param)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Ok(c, result)
}
