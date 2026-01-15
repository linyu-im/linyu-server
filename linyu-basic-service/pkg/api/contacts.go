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
	route.Register("POST", "/basic/v1/contacts/list", ContactsListHandler, false)
	route.Register("POST", "/basic/v1/contacts/rel/delete", ContactsRelDelHandler, false)
}

// ContactsListHandler 通讯录列表
func ContactsListHandler(c *gin.Context) {
	currentUserId := c.GetString("userId")
	list, err := basicService.ContactsService.ContactsList(currentUserId)
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
