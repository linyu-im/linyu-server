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
	route.Register("POST", "/basic/v1/contacts/friend/list", ContactsFriendListHandler)
	route.Register("POST", "/basic/v1/contacts/group/list", ContactsGroupListHandler)
	route.Register("POST", "/basic/v1/contacts/enterprise/list", ContactsEnterpriseListHandler)
	route.Register("POST", "/basic/v1/contacts/friend/delete", ContactsRelDelHandler)
	route.Register("POST", "/basic/v1/contacts/friend/is", ContactsIsFriendHandler)
	route.Register("POST", "/basic/v1/contacts/remark/update", ContactsUpdateRemarkHandler)
	route.Register("POST", "/basic/v1/contacts/tag/update", ContactsUpdateTagHandler)
}

// ContactsFriendListHandler 通讯录好友列表
func ContactsFriendListHandler(c *gin.Context) {
	currentUserId := c.GetString("userId")
	list, err := basicService.ContactsService.ContactsFriendList(currentUserId)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Ok(c, list)
}

// ContactsGroupListHandler 通讯录群列表
func ContactsGroupListHandler(c *gin.Context) {
	currentUserId := c.GetString("userId")
	list, err := basicService.ContactsService.ContactsGroupList(currentUserId)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Ok(c, list)
}

// ContactsEnterpriseListHandler 通讯录企业列表
func ContactsEnterpriseListHandler(c *gin.Context) {
	currentUserId := c.GetString("userId")
	list, err := basicService.ContactsService.ContactsEnterpriseList(currentUserId)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Ok(c, list)
}

// ContactsFriendDelHandler 通讯录好友删除
func ContactsRelDelHandler(c *gin.Context) {
	param := &basicParam.ContactsFriendDeleteParam{}
	if !utils.ShouldBindBodyWithJSONAndValidate(c, param) {
		return
	}
	currentUserId := c.GetString("userId")
	err := basicService.ContactsService.ContactsFriendDelete(currentUserId, param)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Ok(c)
}

// ContactsIsFriendHandler 对方是否是好友
func ContactsIsFriendHandler(c *gin.Context) {
	param := &basicParam.ContactsIsFriendParam{}
	if !utils.ShouldBindBodyWithJSONAndValidate(c, param) {
		return
	}
	currentUserId := c.GetString("userId")
	b := basicService.ContactsService.IsFriend(currentUserId, param.UserId)
	response.Ok(c, b)
}

// ContactsUpdateRemarkHandler 修改好友备注
func ContactsUpdateRemarkHandler(c *gin.Context) {
	param := &basicParam.ContactsUpdateRemarkParam{}
	if !utils.ShouldBindBodyWithJSONAndValidate(c, param) {
		return
	}
	currentUserId := c.GetString("userId")
	err := basicService.ContactsService.UpdateRemark(currentUserId, param)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Ok(c)
}

// ContactsUpdateTagHandler 修改好友标签
func ContactsUpdateTagHandler(c *gin.Context) {
	param := &basicParam.ContactsUpdateTagParam{}
	if !utils.ShouldBindBodyWithJSONAndValidate(c, param) {
		return
	}
	currentUserId := c.GetString("userId")
	err := basicService.ContactsService.UpdateTag(currentUserId, param)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Ok(c)
}
