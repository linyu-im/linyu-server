package api

import (
	"fmt"
	"path/filepath"

	"github.com/gin-gonic/gin"
	basicParam "github.com/linyu-im/linyu-server/linyu-basic-service/pkg/param"
	basicService "github.com/linyu-im/linyu-server/linyu-basic-service/pkg/service"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/response"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/route"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/storage"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/utils"
)

func init() {
	route.Register("POST", "/basic/v1/user/info", UserInfoHandler)
	route.Register("POST", "/basic/v1/user/current/info", CurrentUserInfoHandler)
	route.Register("POST", "/basic/v1/user/emotion/set", UserEmotionSetHandler)
	route.Register("POST", "/basic/v1/user/avatar/get", GetUserAvatarHandler)
	route.Register("POST", "/basic/v1/user/avatar/upload", UploadUserAvatarHandler)
	route.Register("POST", "/basic/v1/user/profile/update", UpdateUserProfileHandler)
}

func CurrentUserInfoHandler(c *gin.Context) {
	currentUserId := c.GetString("userId")
	info, err := basicService.UserService.CurrentUserInfoById(currentUserId)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Ok(c, info)
}

func UserInfoHandler(c *gin.Context) {
	param := &basicParam.UserInfoParam{}
	if !utils.ShouldBindBodyWithJSONAndValidate(c, param) {
		return
	}
	currentUserId := c.GetString("userId")
	info, err := basicService.UserService.UserInfoById(param.UserId, currentUserId)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Ok(c, info)
}

func UserEmotionSetHandler(c *gin.Context) {
	param := &basicParam.UserEmotionSetParam{}
	if !utils.ShouldBindBodyWithJSONAndValidate(c, param) {
		return
	}
	currentUserId := c.GetString("userId")
	err := basicService.EmotionService.SetEmotion(currentUserId, param.EmotionId)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Ok(c)
}

func GetUserAvatarHandler(c *gin.Context) {
	param := &basicParam.GetUserAvatarParam{}
	if !utils.ShouldBindBodyWithJSONAndValidate(c, param) {
		return
	}
	url := basicService.UserService.GetAvatar(param.UserId)
	response.Ok(c, url)
}

// UploadUserAvatarHandler 上传用户头像
func UploadUserAvatarHandler(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		response.Fail(c, "param.file-not-found")
		return
	}

	const maxFileSize = 10 * 1024 * 1024
	if file.Size > maxFileSize {
		response.Fail(c, "file too large, max 10MB")
		return
	}

	currentUserId := c.GetString("userId")
	ext := filepath.Ext(file.Filename)
	fileKey := fmt.Sprintf("avatar/%s%s", currentUserId, ext)

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

	if err := basicService.UserService.UpdateAvatar(currentUserId, url); err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Ok(c, url)
}

// UpdateUserProfileHandler 修改用户资料
func UpdateUserProfileHandler(c *gin.Context) {
	param := &basicParam.UserUpdateProfileParam{}
	if !utils.ShouldBindBodyWithJSONAndValidate(c, param) {
		return
	}
	currentUserId := c.GetString("userId")
	err := basicService.UserService.UpdateProfile(currentUserId, param)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Ok(c)
}
