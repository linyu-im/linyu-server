package api

import (
	"github.com/gin-gonic/gin"
	"github.com/linyu-im/linyu-server/linyu-auth/pkg/service"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/response"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/route"
	"net/http"
)

func init() {
	route.RegisterWhite("GET", "/auth/v1/oauth2/url", Oauth2UrlHandler)
}

// Oauth2UrlHandler oauth2授权地址
func Oauth2UrlHandler(c *gin.Context) {
	oauth2Type := c.Query("type")
	authURL, err := service.Oauth2Service.GetAuthURL(oauth2Type)
	if err != nil {
		response.Fail(c, err.Error())
	}
	c.Redirect(http.StatusFound, authURL)
}
