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
	route.Register("POST", "/cloud-drive/v1/space/user/info", UserSpaceInfoHandler)
	route.Register("POST", "/cloud-drive/v1/space/user/file/list", UserSpaceFileListHandler)
	route.Register("POST", "/cloud-drive/v1/space/user/file/list/all", UserSpaceFileListAllHandler)
	route.Register("POST", "/cloud-drive/v1/space/user/dir/create", UserSpaceDirCreateHandler)
	route.Register("POST", "/cloud-drive/v1/space/user/dir/tree", UserSpaceDirTreeHandler)
	route.Register("POST", "/cloud-drive/v1/space/user/file/delete", UserSpaceFileDeleteHandler)
	route.Register("POST", "/cloud-drive/v1/space/user/file/move", UserSpaceFileMoveHandler)
	route.Register("POST", "/cloud-drive/v1/space/user/file/rename", UserSpaceFileRenameHandler)
	route.Register("POST", "/cloud-drive/v1/space/user/file/category/stats", UserSpaceFileCategoryStatsHandler)
}

// UserSpaceInfoHandler 查询当前用户个人空间信息
func UserSpaceInfoHandler(c *gin.Context) {
	currentUserId := c.GetString("userId")
	space, err := driveService.SpaceService.GetOrCreateUserSpace(currentUserId)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Ok(c, space)
}

// UserSpaceFileListHandler 查询当前用户空间指定目录下的列表
func UserSpaceFileListHandler(c *gin.Context) {
	param := &driveParam.SpaceFileListParam{}
	if !utils.ShouldBindBodyWithJSONAndValidate(c, param) {
		return
	}
	currentUserId := c.GetString("userId")
	list, err := driveService.SpaceService.ListUserSpaceFiles(currentUserId, param.ParentID)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Ok(c, list)
}

// UserSpaceFileListAllHandler 查询目录下全部文件（含子目录，仅文件不含目录）
func UserSpaceFileListAllHandler(c *gin.Context) {
	param := &driveParam.SpaceFileListAllParam{}
	if !utils.ShouldBindBodyWithJSONAndValidate(c, param) {
		return
	}
	currentUserId := c.GetString("userId")
	list, err := driveService.SpaceService.ListUserSpaceAllFiles(currentUserId, param.ParentID)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Ok(c, list)
}

// UserSpaceDirCreateHandler 在当前用户空间指定目录下创建文件夹
func UserSpaceDirCreateHandler(c *gin.Context) {
	param := &driveParam.SpaceCreateDirParam{}
	if !utils.ShouldBindBodyWithJSONAndValidate(c, param) {
		return
	}
	currentUserId := c.GetString("userId")
	dir, err := driveService.SpaceService.CreateUserSpaceDir(currentUserId, param.ParentID, param.DirName)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Ok(c, dir)
}

// UserSpaceDirTreeHandler 查询当前用户空间目录树（仅目录）
func UserSpaceDirTreeHandler(c *gin.Context) {
	currentUserId := c.GetString("userId")
	tree, err := driveService.SpaceService.ListUserSpaceDirTree(currentUserId)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Ok(c, tree)
}

// UserSpaceFileDeleteHandler 删除当前用户空间下的文件或目录
func UserSpaceFileDeleteHandler(c *gin.Context) {
	param := &driveParam.SpaceFileDeleteParam{}
	if !utils.ShouldBindBodyWithJSONAndValidate(c, param) {
		return
	}
	currentUserId := c.GetString("userId")
	if err := driveService.SpaceService.DeleteUserSpaceFiles(currentUserId, param.SpaceFileIDs); err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Ok(c, nil)
}

// UserSpaceFileMoveHandler 将文件或目录移动到其他目录
func UserSpaceFileMoveHandler(c *gin.Context) {
	param := &driveParam.SpaceFileMoveParam{}
	if !utils.ShouldBindBodyWithJSONAndValidate(c, param) {
		return
	}
	currentUserId := c.GetString("userId")
	if err := driveService.SpaceService.MoveUserSpaceFiles(currentUserId, param.SpaceFileIDs, param.TargetParentID); err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Ok(c, nil)
}

// UserSpaceFileRenameHandler 重命名文件或目录
func UserSpaceFileRenameHandler(c *gin.Context) {
	param := &driveParam.SpaceFileRenameParam{}
	if !utils.ShouldBindBodyWithJSONAndValidate(c, param) {
		return
	}
	currentUserId := c.GetString("userId")
	file, err := driveService.SpaceService.RenameUserSpaceFile(currentUserId, param.SpaceFileID, param.NewName)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Ok(c, file)
}

// UserSpaceFileCategoryStatsHandler 查询当前用户空间各文件分类的数量与大小
func UserSpaceFileCategoryStatsHandler(c *gin.Context) {
	currentUserId := c.GetString("userId")
	list, err := driveService.SpaceService.ListUserSpaceCategoryStats(currentUserId)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Ok(c, list)
}
