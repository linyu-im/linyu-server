package api

import (
	"github.com/gin-gonic/gin"
	basicService "github.com/linyu-im/linyu-server/linyu-basic-service/pkg/service"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/config"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/response"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/route"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/utils"
	voipParam "github.com/linyu-im/linyu-server/linyu-voip-chat/pkg/param"
	voipService "github.com/linyu-im/linyu-server/linyu-voip-chat/pkg/service"
)

func init() {
	route.Register("POST", "/voip/v1/livekit/host", LivekitHostHandler)
	route.Register("POST", "/voip/v1/livekit/token/group", LivekitTokenHandler)
}

func LivekitHostHandler(c *gin.Context) {
	response.Ok(c, config.C.Livekit.Host)
}

func LivekitTokenHandler(c *gin.Context) {
	param := &voipParam.LivekitTokenParam{}
	if !utils.ShouldBindBodyWithJSONAndValidate(c, param) {
		return
	}
	currentUserId := c.GetString("userId")
	//验证用户是否属于该群聊
	is := basicService.GroupService.IsGroupMember(param.GroupId, currentUserId)
	if !is {
		response.Fail(c, "param.error")
		return
	}
	token, err := voipService.LivekitService.LivekitToken(param.GroupId, currentUserId)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Ok(c, token)
}
