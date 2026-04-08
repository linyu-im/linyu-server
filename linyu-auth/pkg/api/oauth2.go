package api

import (
	"github.com/gin-gonic/gin"
	"github.com/linyu-im/linyu-server/linyu-auth/pkg/param"
	"github.com/linyu-im/linyu-server/linyu-auth/pkg/service"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/response"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/route"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/utils"
)

func init() {
	route.RegisterWhite("POST", "/auth/v1/oauth2/url", Oauth2UrlHandler)
}

// Oauth2UrlHandler oauth2授权地址
func Oauth2UrlHandler(c *gin.Context) {
	oauth2UrlParam := &param.Oauth2UrlParam{}
	if !utils.ShouldBindBodyWithJSONAndValidate(c, oauth2UrlParam) {
		return
	}
	result, err := service.Oauth2Service.GetAuthURL(oauth2UrlParam.Type)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Ok(c, result)
}
