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
	route.Register("POST", "/basic/v1/contacts/rel/delete", ContactsRelDelHandler)
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

// ContactsRelDelHandler 通讯录关系删除
func ContactsRelDelHandler(c *gin.Context) {
	param := &basicParam.ContactsRelDeleteParam{}
	if !utils.ShouldBindBodyWithJSONAndValidate(c, param) {
		return
	}
	currentUserId := c.GetString("userId")
	err := basicService.ContactsService.ContactsRelDelete(currentUserId, param)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Ok(c)
}
