package api

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	basicParam "github.com/linyu-im/linyu-server/linyu-basic-service/pkg/param"
	basicService "github.com/linyu-im/linyu-server/linyu-basic-service/pkg/service"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/response"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/route"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/storage"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/utils"
)

func init() {
	route.Register("POST", "/basic/v1/feedback/image/upload", FeedbackImageUploadHandler)
	route.Register("POST", "/basic/v1/feedback/create", FeedbackCreateHandler)
}

// FeedbackImageUploadHandler 上传反馈图片
func FeedbackImageUploadHandler(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		response.Fail(c, "param.file-not-found")
		return
	}

	const maxFileSize = 50 * 1024 * 1024
	if file.Size > maxFileSize {
		response.Fail(c, "file too large, max 50MB")
		return
	}

	ext := filepath.Ext(file.Filename)
	datePrefix := time.Now().Format("2006-01-02")
	fileKey := fmt.Sprintf("feedback/%s/%s%s", datePrefix, utils.GenerateUuid(), ext)

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
	response.Ok(c, url)
}

// FeedbackCreateHandler 创建反馈
func FeedbackCreateHandler(c *gin.Context) {
	param := &basicParam.FeedbackCreateParam{}
	if !utils.ShouldBindBodyWithJSONAndValidate(c, param) {
		return
	}
	currentUserId := c.GetString("userId")
	err := basicService.FeedbackService.CreateFeedback(currentUserId, param)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Ok(c)
}
