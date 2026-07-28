package api

import (
	"github.com/gin-gonic/gin"
	driveParam "github.com/linyu-im/linyu-server/linyu-cloud-drive/pkg/param"
	driveService "github.com/linyu-im/linyu-server/linyu-cloud-drive/pkg/service"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/response"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/route"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/utils"
)

func init() {
	route.Register("POST", "/cloud-drive/v1/space-recycle/user/list", UserSpaceRecycleListHandler)
	route.Register("POST", "/cloud-drive/v1/space-recycle/user/restore", UserSpaceRecycleRestoreHandler)
	route.Register("POST", "/cloud-drive/v1/space-recycle/user/delete", UserSpaceRecycleDeleteHandler)
	route.Register("POST", "/cloud-drive/v1/space-recycle/user/clear", UserSpaceRecycleClearHandler)
}

// UserSpaceRecycleListHandler 查询当前用户回收站列表
func UserSpaceRecycleListHandler(c *gin.Context) {
	currentUserId := c.GetString("userId")
	list, err := driveService.SpaceRecycleService.ListUserRecycle(currentUserId)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Ok(c, list)
}

// UserSpaceRecycleRestoreHandler 还原回收站内容
func UserSpaceRecycleRestoreHandler(c *gin.Context) {
	param := &driveParam.SpaceRecycleRestoreParam{}
	if !utils.ShouldBindBodyWithJSONAndValidate(c, param) {
		return
	}
	currentUserId := c.GetString("userId")
	if err := driveService.SpaceRecycleService.RestoreUserRecycle(currentUserId, param.SpaceRecycleIDs); err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Ok(c, nil)
}

// UserSpaceRecycleDeleteHandler 彻底删除回收站内容
func UserSpaceRecycleDeleteHandler(c *gin.Context) {
	param := &driveParam.SpaceRecycleDeleteParam{}
	if !utils.ShouldBindBodyWithJSONAndValidate(c, param) {
		return
	}
	currentUserId := c.GetString("userId")
	if err := driveService.SpaceRecycleService.PermanentlyDeleteUserRecycle(currentUserId, param.SpaceRecycleIDs); err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Ok(c, nil)
}

// UserSpaceRecycleClearHandler 清空当前用户回收站
func UserSpaceRecycleClearHandler(c *gin.Context) {
	currentUserId := c.GetString("userId")
	if err := driveService.SpaceRecycleService.ClearUserRecycle(currentUserId); err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Ok(c, nil)
}
