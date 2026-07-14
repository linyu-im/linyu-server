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
	route.Register("POST", "/basic/v1/group/info/update", GroupUpdateInfoHandler)
	route.Register("POST", "/basic/v1/group/is-admin", GroupIsAdminHandler)
	route.Register("POST", "/basic/v1/group/set-admin", GroupSetAdminHandler)
	route.Register("POST", "/basic/v1/group/transfer-owner", GroupTransferOwnerHandler)
	route.Register("POST", "/basic/v1/group/leave", GroupLeaveHandler)
	route.Register("POST", "/basic/v1/group/notice/list", GroupNoticeListHandler)
	route.Register("POST", "/basic/v1/group/notice/add", GroupNoticeAddHandler)
	route.Register("POST", "/basic/v1/group/notice/update", GroupNoticeUpdateHandler)
	route.Register("POST", "/basic/v1/group/notice/delete", GroupNoticeDeleteHandler)
	route.Register("POST", "/basic/v1/group/nickname/update", GroupUpdateNickNameHandler)
	route.Register("POST", "/basic/v1/group/member/info", GroupMemberInfoHandler)
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

// GroupIsAdminHandler 当前用户是否是群管理员
func GroupIsAdminHandler(c *gin.Context) {
	param := &basicParam.GroupInfoParam{}
	if !utils.ShouldBindBodyWithJSONAndValidate(c, param) {
		return
	}
	currentUserId := c.GetString("userId")
	isAdmin := basicService.GroupService.IsGroupAdmin(param.GroupId, currentUserId)
	response.Ok(c, isAdmin)
}

// GroupUpdateInfoHandler 修改群聊信息
func GroupUpdateInfoHandler(c *gin.Context) {
	param := &basicParam.GroupUpdateInfoParam{}
	if !utils.ShouldBindBodyWithJSONAndValidate(c, param) {
		return
	}
	currentUserId := c.GetString("userId")
	err := basicService.GroupService.UpdateInfo(currentUserId, param)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Ok(c)
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

// GroupTransferOwnerHandler 转让群拥有者
func GroupTransferOwnerHandler(c *gin.Context) {
	param := &basicParam.GroupTransferOwnerParam{}
	if !utils.ShouldBindBodyWithJSONAndValidate(c, param) {
		return
	}
	currentUserId := c.GetString("userId")
	err := basicService.GroupService.TransferOwner(currentUserId, param)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Ok(c)
}

// GroupSetAdminHandler 设置群管理员
func GroupSetAdminHandler(c *gin.Context) {
	param := &basicParam.GroupSetAdminParam{}
	if !utils.ShouldBindBodyWithJSONAndValidate(c, param) {
		return
	}
	currentUserId := c.GetString("userId")
	err := basicService.GroupService.SetAdmin(currentUserId, param)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Ok(c)
}

// GroupLeaveHandler 退出群聊
func GroupLeaveHandler(c *gin.Context) {
	param := &basicParam.GroupLeaveParam{}
	if !utils.ShouldBindBodyWithJSONAndValidate(c, param) {
		return
	}
	currentUserId := c.GetString("userId")
	err := basicService.GroupService.LeaveGroup(currentUserId, param)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Ok(c)
}

// GroupNoticeListHandler 群公告列表
func GroupNoticeListHandler(c *gin.Context) {
	param := &basicParam.GroupNoticeListParam{}
	if !utils.ShouldBindBodyWithJSONAndValidate(c, param) {
		return
	}
	currentUserId := c.GetString("userId")
	data, err := basicService.GroupService.NoticeList(currentUserId, param)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Ok(c, data)
}

// GroupNoticeAddHandler 新增群公告
func GroupNoticeAddHandler(c *gin.Context) {
	param := &basicParam.GroupNoticeAddParam{}
	if !utils.ShouldBindBodyWithJSONAndValidate(c, param) {
		return
	}
	currentUserId := c.GetString("userId")
	err := basicService.GroupService.NoticeAdd(currentUserId, param)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Ok(c)
}

// GroupNoticeUpdateHandler 更新群公告
func GroupNoticeUpdateHandler(c *gin.Context) {
	param := &basicParam.GroupNoticeUpdateParam{}
	if !utils.ShouldBindBodyWithJSONAndValidate(c, param) {
		return
	}
	currentUserId := c.GetString("userId")
	err := basicService.GroupService.NoticeUpdate(currentUserId, param)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Ok(c)
}

// GroupNoticeDeleteHandler 删除群公告
func GroupNoticeDeleteHandler(c *gin.Context) {
	param := &basicParam.GroupNoticeDeleteParam{}
	if !utils.ShouldBindBodyWithJSONAndValidate(c, param) {
		return
	}
	currentUserId := c.GetString("userId")
	err := basicService.GroupService.NoticeDelete(currentUserId, param)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Ok(c)
}

// GroupUpdateNickNameHandler 修改群昵称
func GroupUpdateNickNameHandler(c *gin.Context) {
	param := &basicParam.GroupUpdateNickNameParam{}
	if !utils.ShouldBindBodyWithJSONAndValidate(c, param) {
		return
	}
	currentUserId := c.GetString("userId")
	err := basicService.GroupService.UpdateNickName(currentUserId, param)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Ok(c)
}

// GroupMemberInfoHandler 获取指定成员信息
func GroupMemberInfoHandler(c *gin.Context) {
	param := &basicParam.GroupMemberInfoParam{}
	if !utils.ShouldBindBodyWithJSONAndValidate(c, param) {
		return
	}
	currentUserId := c.GetString("userId")
	member, err := basicService.GroupService.GroupMemberInfo(currentUserId, param)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Ok(c, member)
}
