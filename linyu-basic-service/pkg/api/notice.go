package api

import (
	"github.com/gin-gonic/gin"
	basicService "github.com/linyu-im/linyu-server/linyu-basic-service/pkg/service"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/response"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/route"
)

func init() {
	route.Register("POST", "/basic/v1/notice/list/group", NoticeGroupListHandler)
}

// NoticeGroupListHandler 查询群通知列表
func NoticeGroupListHandler(c *gin.Context) {
	currentUserId := c.GetString("userId")
	list, err := basicService.NoticeService.GroupList(currentUserId)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Ok(c, list)
}
