package api

import (
	"github.com/gin-gonic/gin"
	basicService "github.com/linyu-im/linyu-server/linyu-basic-service/pkg/service"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/response"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/route"
)

func init() {
	route.Register("POST", "/basic/v1/user-badge/list", UserBadgeListHandler)
}

// UserBadgeListHandler 查询当前用户红点数据
func UserBadgeListHandler(c *gin.Context) {
	currentUserId := c.GetString("userId")
	list, err := basicService.UserBadgeService.List(currentUserId)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Ok(c, list)
}
