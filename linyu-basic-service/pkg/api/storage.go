package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	basicParam "github.com/linyu-im/linyu-server/linyu-basic-service/pkg/param"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/config"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/response"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/route"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/storage"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/utils"
	"net/http"
	"path/filepath"
	"strings"
	"time"
)

func init() {
	route.RegisterWhite("GET", storage.LocalStorageUrl+"/*fileKey", LocalStorageHandler)
	route.Register("POST", "/basic/v1/storage/upload", StorageUploadHandler)
	route.Register("POST", "/basic/v1/storage/delete", StorageDeleteHandler)
}

// StorageDeleteHandler 文件删除
func StorageDeleteHandler(c *gin.Context) {
	param := &basicParam.StorageDeleteParam{}
	if !utils.ShouldBindBodyWithJSONAndValidate(c, param) {
		return
	}
	currentUserId := c.GetString("userId")

	parts := strings.Split(strings.TrimPrefix(param.FileKey, "/"), "/")
	if parts[1] != currentUserId {
		response.Fail(c, "param.error")
		return
	}

	err := storage.S.Delete(param.FileKey)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}

	response.Ok(c)
}

// StorageUploadHandler 文件上传
func StorageUploadHandler(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		response.Fail(c, "file-not-found")
		return
	}
	currentUserId := c.GetString("userId")
	ext := filepath.Ext(file.Filename)
	datePrefix := time.Now().Format("2006-01-02")
	fileKey := fmt.Sprintf("%s/%s/%s%s", datePrefix, currentUserId, utils.GenerateUuid(), ext)

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

// LocalStorageHandler 本地文件内容获取
func LocalStorageHandler(c *gin.Context) {
	fileKey := c.Param("fileKey")
	fileKey = strings.TrimPrefix(fileKey, "/")

	if fileKey == "" {
		c.Status(http.StatusNotFound)
		return
	}
	if config.C.Storage.Type != config.LocalStorageType {
		response.Fail(c, "basic.storage.local-storage-not-exist")
		return
	}

	localStorage, ok := storage.S.(*storage.LocalStorage)
	if !ok {
		c.Status(http.StatusInternalServerError)
		return
	}

	absPath, err := utils.GetFileSystemPath(localStorage.FilePath, fileKey)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}

	c.File(absPath)
}
