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
	route.Register("POST", "/basic/v1/user/current/info", CurrentUserInfoHandler)
	route.Register("POST", "/basic/v1/user/emotion/set", UserEmotionSetHandler)
	route.Register("POST", "/basic/v1/user/avatar/get", GetUserAvatarHandler)
}

func CurrentUserInfoHandler(c *gin.Context) {
	currentUserId := c.GetString("userId")
	info, err := basicService.UserService.UserInfoById(currentUserId)
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
