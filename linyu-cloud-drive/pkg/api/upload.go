package api

import (
	"github.com/gin-gonic/gin"
	driveParam "github.com/linyu-im/linyu-server/linyu-cloud-drive/pkg/param"
	driveResult "github.com/linyu-im/linyu-server/linyu-cloud-drive/pkg/result"
	driveService "github.com/linyu-im/linyu-server/linyu-cloud-drive/pkg/service"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/response"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/route"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/storage"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/utils"
)

func init() {
	route.Register("POST", "/cloud-drive/v1/upload/check", CheckUploadStatusHandler)
	route.Register("POST", "/cloud-drive/v1/upload/chunk", UploadChunkHandler)
	route.Register("POST", "/cloud-drive/v1/upload/merge", MergeChunkHandler)

}

// CheckUploadStatusHandler 校验文件上传状态
func CheckUploadStatusHandler(c *gin.Context) {
	param := &driveParam.UploadFileInfoParam{}
	if !utils.ShouldBindBodyWithJSONAndValidate(c, param) {
		return
	}
	currentUserId := c.GetString("userId")

	// 验证物理文件是否存在(秒传)
	file := driveService.PhysicalFileService.GetFileByHash(param.FileHash)
	if file != nil {
		err := driveService.SpaceFileService.CreateFileFromPhysicalFile(currentUserId, param, file)
		if err != nil {
			response.Fail(c, err.Error())
			return
		}
		response.Ok(c, &driveResult.CheckUploadStatusResult{
			Uploaded: true,
		})
		return
	}

	response.Ok(c, &driveResult.CheckUploadStatusResult{
		Uploaded: false,
		// 获取已上传分片
		UploadedChunks: storage.GetUploadChunkInfo(param.FileHash),
	})
}

// UploadChunkHandler 上传切片
func UploadChunkHandler(c *gin.Context) {
	fileHash := c.PostForm("fileHash")
	chunkIndex := c.PostForm("chunkIndex")
	if fileHash == "" || chunkIndex == "" {
		response.Fail(c, "param.error")
		return
	}
	file, err := c.FormFile("file")
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	if err = storage.UploadChunk(file, chunkIndex, fileHash); err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Ok(c)
}

// MergeChunkHandler 合并切片
func MergeChunkHandler(c *gin.Context) {
	param := &driveParam.UploadFileInfoParam{}
	if !utils.ShouldBindBodyWithJSONAndValidate(c, param) {
		return
	}
	currentUserId := c.GetString("userId")
	err, storagePath := storage.MergeChunk(param.FileHash, param.TotalChunk, param.FileName, "clouddrive")

	err, physicalFile := driveService.PhysicalFileService.CreateFile(param, storagePath)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	err = driveService.SpaceFileService.CreateFileFromPhysicalFile(currentUserId, param, physicalFile)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Ok(c, storagePath)
}
